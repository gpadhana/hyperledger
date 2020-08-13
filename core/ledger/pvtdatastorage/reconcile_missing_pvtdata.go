/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package pvtdatastorage

import (
	"github.com/hyperledger/fabric-protos-go/ledger/rwset"
	"github.com/hyperledger/fabric/common/ledger/util/leveldbhelper"
	"github.com/hyperledger/fabric/core/ledger"
	"github.com/pkg/errors"
	"github.com/willf/bitset"
)

// CommitPvtDataOfOldBlocks commits the pvtData (i.e., previously missing data) of old blockp.
// The parameter `blocksPvtData` refers a list of old block's pvtdata which are missing in the pvtstore.
// Given a list of old block's pvtData, `CommitPvtDataOfOldBlocks` performs the following three
// operations
// (1) construct update entries (i.e., dataEntries, expiryEntries, missingDataEntries)
//     from the above created data entries
// (2) create a db update batch from the update entries
// (3) commit the update batch to the pvtStore
func (s *Store) CommitPvtDataOfOldBlocks(
	blocksPvtData map[uint64][]*ledger.TxPvtData,
	deprioritizedMissingData ledger.MissingPvtDataInfo,
) error {
	s.purgerLock.Lock()
	defer s.purgerLock.Unlock()

	if s.isLastUpdatedOldBlocksSet {
		return &ErrIllegalCall{`The lastUpdatedOldBlocksList is set. It means that the
		stateDB may not be in sync with the pvtStore`}
	}

	p := &oldBlockDataProcessor{
		Store: s,
		entries: &entriesForPvtDataOfOldBlocks{
			dataEntries:                     make(map[dataKey]*rwset.CollectionPvtReadWriteSet),
			expiryEntries:                   make(map[expiryKey]*ExpiryData),
			prioritizedMissingDataEntries:   make(map[nsCollBlk]*bitset.BitSet),
			deprioritizedMissingDataEntries: make(map[nsCollBlk]*bitset.BitSet),
		},
	}

	if err := p.prepareDataAndExpiryEntries(blocksPvtData); err != nil {
		return err
	}

	if err := p.prepareMissingDataEntriesToReflectReconciledData(); err != nil {
		return err
	}

	if err := p.prepareMissingDataEntriesToReflectPriority(deprioritizedMissingData); err != nil {
		return err
	}

	batch, err := p.constructDBUpdateBatch()
	if err != nil {
		return err
	}
	return s.db.WriteBatch(batch, true)
}

type oldBlockDataProcessor struct {
	*Store
	entries *entriesForPvtDataOfOldBlocks
}

func (p *oldBlockDataProcessor) prepareDataAndExpiryEntries(blocksPvtData map[uint64][]*ledger.TxPvtData) error {
	var dataEntries []*dataEntry
	var expData *ExpiryData

	for blkNum, pvtData := range blocksPvtData {
		dataEntries = append(dataEntries, prepareDataEntries(blkNum, pvtData)...)
	}

	for _, dataEntry := range dataEntries {
		nsCollBlk := dataEntry.key.nsCollBlk
		txNum := dataEntry.key.txNum

		expKey, err := p.constructExpiryKey(dataEntry)
		if err != nil {
			return err
		}

		if neverExpires(expKey.expiringBlk) {
			p.entries.dataEntries[*dataEntry.key] = dataEntry.value
			continue
		}

		if expData, err = p.getExpiryData(expKey); err != nil {
			return err
		}
		if expData == nil {
			// if expiryData is not available, it means that
			// the pruge scheduler removed these entries and the
			// associated data entry is no longer needed. Note
			// that the associated missingData entry would also
			// be not present. Hence, we can skip this data entry.
			continue
		}
		expData.addPresentData(nsCollBlk.ns, nsCollBlk.coll, txNum)

		p.entries.dataEntries[*dataEntry.key] = dataEntry.value
		p.entries.expiryEntries[expKey] = expData
	}
	return nil
}

func (p *oldBlockDataProcessor) prepareMissingDataEntriesToReflectReconciledData() error {
	for dataKey := range p.entries.dataEntries {
		key := dataKey.nsCollBlk
		txNum := uint(dataKey.txNum)

		prioMissingData, err := p.getMissingDataFromPrioritizedList(key)
		if err != nil {
			return err
		}
		if prioMissingData != nil && prioMissingData.Test(txNum) {
			p.entries.prioritizedMissingDataEntries[key] = prioMissingData.Clear(txNum)
			continue
		}

		deprioMissingData, err := p.getMissingDataFromDeprioritizedList(key)
		if err != nil {
			return err
		}
		if deprioMissingData != nil && deprioMissingData.Test(txNum) {
			p.entries.deprioritizedMissingDataEntries[key] = deprioMissingData.Clear(txNum)
		}
		// if the missing data entry is already purged by the purge scheduler, we would
		// get nil missingData from both prioritized and deprioritized list
	}

	return nil
}

func (p *oldBlockDataProcessor) prepareMissingDataEntriesToReflectPriority(deprioritizedList ledger.MissingPvtDataInfo) error {
	for blkNum, blkMissingData := range deprioritizedList {
		for txNum, txMissingData := range blkMissingData {
			for _, nsColl := range txMissingData {
				key := nsCollBlk{
					ns:     nsColl.Namespace,
					coll:   nsColl.Collection,
					blkNum: blkNum,
				}
				txNum := uint(txNum)

				prioMissingData, err := p.getMissingDataFromPrioritizedList(key)
				if err != nil {
					return err
				}
				if prioMissingData == nil {
					// we would reach here when either of the following happens:
					//   (1) when the purge scheduler already removed the respective
					//       missing data entry.
					//   (2) when the missing data info is already persistent in the
					//       deprioritized list. Currently, we do not have different
					//       levels of deprioritized list.
					// In both of the above case, we can continue to the next entry.
					continue
				}
				p.entries.prioritizedMissingDataEntries[key] = prioMissingData.Clear(txNum)

				deprioMissingData, err := p.getMissingDataFromDeprioritizedList(key)
				if err != nil {
					return err
				}
				if deprioMissingData == nil {
					deprioMissingData = &bitset.BitSet{}
				}
				p.entries.deprioritizedMissingDataEntries[key] = deprioMissingData.Set(txNum)
			}
		}
	}

	return nil
}

func (p *oldBlockDataProcessor) constructExpiryKey(dataEntry *dataEntry) (expiryKey, error) {
	// get the expiryBlk number to construct the expiryKey
	nsCollBlk := dataEntry.key.nsCollBlk
	expiringBlk, err := p.btlPolicy.GetExpiringBlock(nsCollBlk.ns, nsCollBlk.coll, nsCollBlk.blkNum)
	if err != nil {
		return expiryKey{}, errors.WithMessagef(err, "error while constructing expiry data key")
	}

	return expiryKey{
		expiringBlk:   expiringBlk,
		committingBlk: nsCollBlk.blkNum,
	}, nil
}

func (p *oldBlockDataProcessor) getExpiryData(expKey expiryKey) (*ExpiryData, error) {
	if expiryData, ok := p.entries.expiryEntries[expKey]; ok {
		return expiryData, nil
	}

	expData, err := p.db.Get(encodeExpiryKey(&expKey))
	if err != nil {
		return nil, err
	}
	if expData == nil {
		return nil, errors.Wrap(err, "error while getting expiry data from the store")
	}
	return decodeExpiryValue(expData)
}

func (p *oldBlockDataProcessor) getMissingDataFromPrioritizedList(nsCollBlk nsCollBlk) (*bitset.BitSet, error) {
	missingData, ok := p.entries.prioritizedMissingDataEntries[nsCollBlk]
	if ok {
		return missingData, nil
	}

	return p.getMissingDataBitmapFromStore(elgPrioritizedMissingDataGroup, nsCollBlk)
}

func (p *oldBlockDataProcessor) getMissingDataFromDeprioritizedList(nsCollBlk nsCollBlk) (*bitset.BitSet, error) {
	missingData, ok := p.entries.deprioritizedMissingDataEntries[nsCollBlk]
	if ok {
		return missingData, nil
	}

	return p.getMissingDataBitmapFromStore(elgDeprioritizedMissingDataGroup, nsCollBlk)
}

func (p *oldBlockDataProcessor) getMissingDataBitmapFromStore(group []byte, nsCollBlk nsCollBlk) (*bitset.BitSet, error) {
	key := encodeElgMissingDataKey(
		group,
		&missingDataKey{
			nsCollBlk: nsCollBlk,
		},
	)

	missingData, err := p.db.Get(key)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting missing data bitmap from the store")
	}
	if missingData == nil {
		return nil, nil
	}

	return decodeMissingDataValue(missingData)
}

func (p *oldBlockDataProcessor) constructDBUpdateBatch() (*leveldbhelper.UpdateBatch, error) {
	batch := p.db.NewUpdateBatch()

	if err := p.entries.addDataEntriesTo(batch); err != nil {
		return nil, errors.WithMessage(err, "error while adding data entries to the update batch")
	}

	if err := p.entries.addExpiryEntriesTo(batch); err != nil {
		return nil, errors.WithMessage(err, "error while adding expiry entries to the update batch")
	}

	if err := p.entries.addMissingDataEntriesTo(batch); err != nil {
		return nil, errors.WithMessage(err, "error while adding missing data entries to the update batch")
	}

	return batch, nil
}

type entriesForPvtDataOfOldBlocks struct {
	dataEntries                     map[dataKey]*rwset.CollectionPvtReadWriteSet
	expiryEntries                   map[expiryKey]*ExpiryData
	prioritizedMissingDataEntries   map[nsCollBlk]*bitset.BitSet
	deprioritizedMissingDataEntries map[nsCollBlk]*bitset.BitSet
}

func (e *entriesForPvtDataOfOldBlocks) addDataEntriesTo(batch *leveldbhelper.UpdateBatch) error {
	var key, val []byte
	var err error

	for dataKey, pvtData := range e.dataEntries {
		key = encodeDataKey(&dataKey)
		if val, err = encodeDataValue(pvtData); err != nil {
			return errors.Wrap(err, "error while encoding data value")
		}
		batch.Put(key, val)
	}
	return nil
}

func (e *entriesForPvtDataOfOldBlocks) addExpiryEntriesTo(batch *leveldbhelper.UpdateBatch) error {
	var key, val []byte
	var err error

	for expiryKey, expiryData := range e.expiryEntries {
		key = encodeExpiryKey(&expiryKey)
		if val, err = encodeExpiryValue(expiryData); err != nil {
			return errors.Wrap(err, "error while encoding expiry value")
		}
		batch.Put(key, val)
	}
	return nil
}

func (e *entriesForPvtDataOfOldBlocks) addMissingDataEntriesTo(batch *leveldbhelper.UpdateBatch) error {
	var key, val []byte
	var err error

	entries := map[string]map[nsCollBlk]*bitset.BitSet{
		string(elgPrioritizedMissingDataGroup):   e.prioritizedMissingDataEntries,
		string(elgDeprioritizedMissingDataGroup): e.deprioritizedMissingDataEntries,
	}

	for group, missingDataList := range entries {
		for nsCollBlk, missingData := range missingDataList {
			key = encodeElgMissingDataKey(
				[]byte(group),
				&missingDataKey{
					nsCollBlk:  nsCollBlk,
					isEligible: true,
				},
			)

			if missingData.None() {
				batch.Delete(key)
				continue
			}

			if val, err = encodeMissingDataValue(missingData); err != nil {
				return errors.Wrap(err, "error while encoding missing data bitmap")
			}
			batch.Put(key, val)
		}
	}
	return nil
}
