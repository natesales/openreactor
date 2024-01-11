package main

type TCP015Controller struct {
	Controller TurboController
}

// Off turns the pump off
func (t *TCP015Controller) Off() error {
	return t.Controller.SetRegister(10, true)
}

// On turns the pump on
func (t *TCP015Controller) On() error {
	return t.Controller.SetRegister(10, false)
}

// Standby puts the pump in standby mode
func (t *TCP015Controller) Standby() error {
	// TODO
	panic("not implemented")
}

// IsRunning returns true if the pump is running
func (t *TCP015Controller) IsRunning() (bool, error) {
	message, err := t.Controller.ReadRegister(10)
	if err != nil {
		return false, err
	}
	return message.Payload == "111111", nil
}

// Hz returns the current motor speed
func (t *TCP015Controller) Hz() (int, error) {
	message, err := t.Controller.ReadRegister(309)
	if err != nil {
		return 0, err
	}
	return toInt(message.Payload), nil
}

// CurrentDraw returns the current motor current draw
func (t *TCP015Controller) CurrentDraw() (int, error) {
	message, err := t.Controller.ReadRegister(310)
	if err != nil {
		return 0, err
	}
	return toInt(message.Payload) / 100, nil
}

// FirmwareVersion returns the firmware version
func (t *TCP015Controller) FirmwareVersion() (string, error) {
	message, err := t.Controller.ReadRegister(312)
	if err != nil {
		return "", err
	}
	return message.Payload, nil
}

// ErrorCode returns the current error code
func (t *TCP015Controller) ErrorCode() (string, error) {
	message, err := t.Controller.ReadRegister(303)
	if err != nil {
		return "", err
	}
	return message.Payload, nil
}
