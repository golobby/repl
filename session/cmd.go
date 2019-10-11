package session

func (s *Session) handleShellCommands(code string) error {
	typ, data := ParseCmd(code)
	switch typ {
	case REPLCmdDoc:
		output, err := goDoc(data)
		if err != nil {
			return err
		}
		s.shellCmdOutput = string(output)
		return nil
	case REPLCmdHelp:
		s.shellCmdOutput = helpText
		return nil
	case REPLCmdTypeVal:
		return s.Add(wrapInPrint(data))
	case REPLCmdPop:
		s.code = s.code[:len(s.code)-1]
		return nil
	case REPLCmdDump:
		s.shellCmdOutput = s.dump()
		return nil
	default:
		return nil
	}
	return nil
}
