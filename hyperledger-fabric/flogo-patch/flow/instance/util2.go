package instance

// Host returns parent flow instanance
func (ti *TaskInst) Host() *TaskInst {
	h := ti.flowInst.host
	ti.logger.Debugf("got flow host: %+v", h)
	if inst, ok := h.(*TaskInst); ok {
		return inst
	}
	return nil
}
