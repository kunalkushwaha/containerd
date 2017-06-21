package metadata

type metadataExistsErr struct {
	desc string
}
type metadataNotFoundErr struct {
	desc string
}
type metadataNotEmptyErr struct {
	desc string
}

// ErrExists is returned when an item already exists in metadata
func ErrExists(msg string) error {
	if msg == "" {
		msg = "metadata: exists"
	}
	return metadataExistsErr{
		desc: msg,
	}
}

// ErrNotFound is returned when an item cannot be found in metadata
func ErrNotFound(msg string) error {
	if msg == "" {
		msg = "metadata: not found"
	}
	return metadataNotFoundErr{
		desc: msg,
	}
}

// ErrNotEmpty is returned when a metadata item can't be deleted because it is not empty
func ErrNotEmpty(msg string) error {
	if msg == "" {
		msg = "metadata: namespace not empty"
	}
	return metadataNotEmptyErr{
		desc: msg,
	}
}

func (m metadataExistsErr) Error() string {
	return m.desc
}
func (m metadataNotFoundErr) Error() string {
	return m.desc
}
func (m metadataNotEmptyErr) Error() string {
	return m.desc
}

func (m metadataExistsErr) Exists() bool {
	return true
}

func (m metadataNotFoundErr) NotFound() bool {
	return true
}

func (m metadataNotEmptyErr) NotEmpty() bool {
	return true
}

// IsNotFound returns true if the error is due to a missing metadata item
func IsNotFound(err error) bool {
	if err, ok := err.(interface {
		NotFound() bool
	}); ok {
		return err.NotFound()
	}

	causal, ok := err.(interface {
		Cause() error
	})
	if !ok {
		return false
	}

	return IsNotFound(causal.Cause())
}

// IsExists returns true if the error is due to an already existing metadata item
func IsExists(err error) bool {
	if err, ok := err.(interface {
		Exists() bool
	}); ok {
		return err.Exists()
	}

	causal, ok := err.(interface {
		Cause() error
	})
	if !ok {
		return false
	}

	return IsExists(causal.Cause())
}

// IsNotEmpty returns true if the error is due to delete request of a non-empty metadata item
func IsNotEmpty(err error) bool {
	if err, ok := err.(interface {
		NotEmpty() bool
	}); ok {
		return err.NotEmpty()
	}

	causal, ok := err.(interface {
		Cause() error
	})
	if !ok {
		return false
	}

	return IsNotEmpty(causal.Cause())
}
