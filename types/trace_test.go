package types

// func TestSender(t *testing.T) {
// 	sender := Sender{}

// 	// Test SetUserID
// 	userID := "abc"
// 	sender.SetUserID(userID)
// 	if sender.GetUserID() != userID {
// 		t.Errorf("Sender: expected set Nonce to %q but got %q", userID, sender.GetUserID())
// 	}

// 	// Test SetUserID
// 	privateKeyID := "def"
// 	sender.SetPrivateKeyID(privateKeyID)
// 	if sender.GetPrivateKeyID() != privateKeyID {
// 		t.Errorf("Sender: expected set Nonce to %q but got %q", privateKeyID, sender.GetPrivateKeyID())
// 	}
// }

// func TestChain(t *testing.T) {
// 	chain := Chain{}

// 	// Test SetUserID
// 	ID := "abc"
// 	chain.SetID(ID)
// 	if chain.GetID() != ID {
// 		t.Errorf("Chain: expected set ID to %q but got %q", ID, chain.GetID())
// 	}

// 	// Test SetUserID
// 	isEIP155 := true
// 	chain.SetEIP155(isEIP155)
// 	if chain.GetEIP155() != isEIP155 {
// 		t.Errorf("Chain: expected set IsEIP155 to %v but got %v", isEIP155, chain.GetEIP155())
// 	}
// }

// func TestReceiver(t *testing.T) {
// 	receiver := Receiver{}

// 	// Test SetUserID
// 	ID := "abc"
// 	receiver.SetID(ID)
// 	if receiver.GetID() != ID {
// 		t.Errorf("Receiver: expected set ID to %q but got %q", ID, receiver.GetID())
// 	}

// 	// Test SetUserID
// 	address := common.HexToAddress("0xfF778b716FC07D98839f48DdB88D8bE583BEB684")
// 	receiver.SetAddress(&address)
// 	if receiver.GetAddress().Hex() != address.Hex() {
// 		t.Errorf("Receiver: expected set IsEIP155 to %v but got %v", address.Hex(), receiver.GetAddress().Hex())
// 	}
// }

// func TestCall(t *testing.T) {
// 	call := Call{}

// 	// Test SetUserID
// 	methodID := "abc"
// 	call.SetMethodID(methodID)
// 	if call.GetMethodID() != methodID {
// 		t.Errorf("Call: expected set MethodID to %q but got %q", methodID, call.GetMethodID())
// 	}

// 	// Test SetUserID
// 	v := hexutil.MustDecodeBig("0x12ae30fd5c9")
// 	call.SetValue(v)
// 	if call.GetValue() != v {
// 		t.Errorf("Call: expected set Value to %q but got %q", v, call.GetValue())
// 	}

// 	args := []string{"0xfF778b716FC07D98839f48DdB88D8bE583BEB684", "0x1"}
// 	call.SetArgs(args)
// 	for i, arg := range call.GetArgs() {
// 		if arg != args[i] {
// 			t.Errorf("Call: expected Arg to be %v but got %v", args[i], arg)
// 		}
// 	}
// }

// func TestTrace(t *testing.T) {
// 	// Init trace object
// 	trace := Trace{}

// 	var (
// 		userID       = "abc"
// 		privateKeyID = "def"
// 	)

// 	sender := Sender{userID, privateKeyID}
// 	trace.SetSender(&sender)

// 	if trace.GetSender().GetUserID() != userID {
// 		t.Errorf("Trace: expected set UserID to %q but got %q", userID, trace.GetSender().GetUserID())
// 	}

// 	var (
// 		chainID  = "abc"
// 		isEIP155 = true
// 	)

// 	chain := Chain{chainID, isEIP155}
// 	trace.SetChain(&chain)

// 	if trace.GetChain().GetID() != chainID {
// 		t.Errorf("Trace: expected set ChainID to %q but got %q", chainID, trace.GetChain().GetID())
// 	}

// 	var (
// 		receiverID = "abc"
// 		address    = common.HexToAddress("0xfF778b716FC07D98839f48DdB88D8bE583BEB684")
// 	)

// 	receiver := Receiver{receiverID, &address}
// 	trace.SetReceiver(&receiver)

// 	if trace.GetReceiver().GetID() != receiverID {
// 		t.Errorf("Trace: expected set ReceiverID to %q but got %q", receiverID, trace.GetReceiver().GetID())
// 	}

// 	var (
// 		methodID = "abc"
// 		value    = hexutil.MustDecodeBig("0x0")
// 		args     = []string{}
// 	)

// 	call := Call{methodID, value, args}
// 	trace.SetCall(&call)

// 	if trace.GetCall().GetMethodID() != methodID {
// 		t.Errorf("Trace: expected set MethodID to %q but got %q", methodID, trace.GetCall().GetMethodID())
// 	}

// 	var (
// 		raw = hexutil.MustDecode("0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80")
// 	)
// 	tx := Transaction{Raw: raw}

// 	trace.SetTx(&tx)
// 	if hexutil.Encode(trace.GetTx().GetRaw()) != hexutil.Encode(raw) {
// 		t.Errorf("Trace: expected set Raw to %q but got %q", hexutil.Encode(raw), hexutil.Encode(trace.GetTx().GetRaw()))
// 	}
// }
