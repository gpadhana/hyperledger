// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ApplicationCapabilities is an autogenerated mock type for the ApplicationCapabilities type
type ApplicationCapabilities struct {
	mock.Mock
}

// ACLs provides a mock function with given fields:
func (_m *ApplicationCapabilities) ACLs() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// CollectionUpgrade provides a mock function with given fields:
func (_m *ApplicationCapabilities) CollectionUpgrade() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ForbidDuplicateTXIdInBlock provides a mock function with given fields:
func (_m *ApplicationCapabilities) ForbidDuplicateTXIdInBlock() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// KeyLevelEndorsement provides a mock function with given fields:
func (_m *ApplicationCapabilities) KeyLevelEndorsement() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// LifecycleV20 provides a mock function with given fields:
func (_m *ApplicationCapabilities) LifecycleV20() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MetadataLifecycle provides a mock function with given fields:
func (_m *ApplicationCapabilities) MetadataLifecycle() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// PrivateChannelData provides a mock function with given fields:
func (_m *ApplicationCapabilities) PrivateChannelData() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// StorePvtDataOfInvalidTx provides a mock function with given fields:
func (_m *ApplicationCapabilities) StorePvtDataOfInvalidTx() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Supported provides a mock function with given fields:
func (_m *ApplicationCapabilities) Supported() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// V1_1Validation provides a mock function with given fields:
func (_m *ApplicationCapabilities) V1_1Validation() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// V1_2Validation provides a mock function with given fields:
func (_m *ApplicationCapabilities) V1_2Validation() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// V1_3Validation provides a mock function with given fields:
func (_m *ApplicationCapabilities) V1_3Validation() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// V20CouchdbValidation provides a mock function with given fields:
func (_m *ApplicationCapabilities) V20CouchdbValidation() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// V2_0Validation provides a mock function with given fields:
func (_m *ApplicationCapabilities) V2_0Validation() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
