package tbsidecar

type AccountFlags struct {
	Linked                     bool
	DebitsMustNotExceedCredits bool
	CreditsMustNotExceedDebits bool
}

func (o Account) AccountFlags() AccountFlags {
	var f AccountFlags
	f.Linked = ((o.Flags >> 0) & 0x1) == 1
	f.DebitsMustNotExceedCredits = ((o.Flags >> 1) & 0x1) == 1
	f.CreditsMustNotExceedDebits = ((o.Flags >> 2) & 0x1) == 1
	return f
}

func (f AccountFlags) ToUint16() uint16 {
	var ret uint16 = 0

	if f.Linked {
		ret |= (1 << 0)
	}

	if f.DebitsMustNotExceedCredits {
		ret |= (1 << 1)
	}

	if f.CreditsMustNotExceedDebits {
		ret |= (1 << 2)
	}

	return ret
}

func (o Transfer) TransferFlags() TransferFlags {
	var f TransferFlags
	f.Linked = ((o.Flags >> 0) & 0x1) == 1
	f.Pending = ((o.Flags >> 1) & 0x1) == 1
	f.PostPendingTransfer = ((o.Flags >> 2) & 0x1) == 1
	f.VoidPendingTransfer = ((o.Flags >> 3) & 0x1) == 1
	f.BalancingDebit = ((o.Flags >> 4) & 0x1) == 1
	f.BalancingCredit = ((o.Flags >> 5) & 0x1) == 1
	return f
}

type TransferFlags struct {
	Linked              bool
	Pending             bool
	PostPendingTransfer bool
	VoidPendingTransfer bool
	BalancingDebit      bool
	BalancingCredit     bool
}

func (f TransferFlags) ToUint16() uint16 {
	var ret uint16 = 0

	if f.Linked {
		ret |= (1 << 0)
	}

	if f.Pending {
		ret |= (1 << 1)
	}

	if f.PostPendingTransfer {
		ret |= (1 << 2)
	}

	if f.VoidPendingTransfer {
		ret |= (1 << 3)
	}

	if f.BalancingDebit {
		ret |= (1 << 4)
	}

	if f.BalancingCredit {
		ret |= (1 << 5)
	}

	return ret
}

type Account struct {
	ID             [16]byte  `json:"id"`
	UserData       [16]byte  `json:"userData"`
	Reserved       [48]uint8 `json:"reserved"`
	Ledger         uint32    `json:"ledger"`
	Code           uint16    `json:"code"`
	Flags          uint16    `json:"flags"`
	DebitsPending  uint64    `json:"debitsPending"`
	DebitsPosted   uint64    `json:"debitsPosted"`
	CreditsPending uint64    `json:"creditsPending"`
	CreditsPosted  uint64    `json:"creditsPosted"`
	Timestamp      uint64    `json:"timestamp"`
}

type Transfer struct {
	ID              [16]byte `json:"id"`
	DebitAccountID  [16]byte `json:"debitAccountId"`
	CreditAccountID [16]byte `json:"creditAccountId"`
	UserData        [16]byte `json:"userData"`
	Reserved        [16]byte `json:"reserved"`
	PendingID       [16]byte `json:"pendingId"`
	Timeout         uint64   `json:"timeout"`
	Ledger          uint32   `json:"ledger"`
	Code            uint16   `json:"code"`
	Flags           uint16   `json:"flags"`
	Amount          uint64   `json:"amount"`
	Timestamp       uint64   `json:"timestamp"`
}

type AccountEventResult struct {
	Index  uint32              `json:"index"`
	Result CreateAccountResult `json:"result"`
}

type TransferEventResult struct {
	Index  uint32               `json:"index"`
	Result CreateTransferResult `json:"result"`
}

type CreateAccountResult uint32

const (
	AccountOK                          CreateAccountResult = 0
	AccountLinkedEventFailed           CreateAccountResult = 1
	AccountLinkedEventChainOpen        CreateAccountResult = 2
	AccountTimestampMustBeZero         CreateAccountResult = 3
	AccountReservedFlag                CreateAccountResult = 4
	AccountReservedField               CreateAccountResult = 5
	AccountIDMustNotBeZero             CreateAccountResult = 6
	AccountIDMustNotBeIntMax           CreateAccountResult = 7
	AccountFlagsAreMutuallyExclusive   CreateAccountResult = 8
	AccountLedgerMustNotBeZero         CreateAccountResult = 9
	AccountCodeMustNotBeZero           CreateAccountResult = 10
	AccountDebitsPendingMustBeZero     CreateAccountResult = 11
	AccountDebitsPostedMustBeZero      CreateAccountResult = 12
	AccountCreditsPendingMustBeZero    CreateAccountResult = 13
	AccountCreditsPostedMustBeZero     CreateAccountResult = 14
	AccountExistsWithDifferentFlags    CreateAccountResult = 15
	AccountExistsWithDifferentUserData CreateAccountResult = 16
	AccountExistsWithDifferentLedger   CreateAccountResult = 17
	AccountExistsWithDifferentCode     CreateAccountResult = 18
	AccountExists                      CreateAccountResult = 19
)

type CreateTransferResult uint32

const (
	TransferOK                                         CreateTransferResult = 0
	TransferLinkedEventFailed                          CreateTransferResult = 1
	TransferLinkedEventChainOpen                       CreateTransferResult = 2
	TransferTimestampMustBeZero                        CreateTransferResult = 3
	TransferReservedFlag                               CreateTransferResult = 4
	TransferReservedField                              CreateTransferResult = 5
	TransferIDMustNotBeZero                            CreateTransferResult = 6
	TransferIDMustNotBeIntMax                          CreateTransferResult = 7
	TransferFlagsAreMutuallyExclusive                  CreateTransferResult = 8
	TransferDebitAccountIDMustNotBeZero                CreateTransferResult = 9
	TransferDebitAccountIDMustNotBeIntMax              CreateTransferResult = 10
	TransferCreditAccountIDMustNotBeZero               CreateTransferResult = 11
	TransferCreditAccountIDMustNotBeIntMax             CreateTransferResult = 12
	TransferAccountsMustBeDifferent                    CreateTransferResult = 13
	TransferPendingIDMustBeZero                        CreateTransferResult = 14
	TransferPendingIDMustNotBeZero                     CreateTransferResult = 15
	TransferPendingIDMustNotBeIntMax                   CreateTransferResult = 16
	TransferPendingIDMustBeDifferent                   CreateTransferResult = 17
	TransferTimeoutReservedForPendingTransfer          CreateTransferResult = 18
	TransferLedgerMustNotBeZero                        CreateTransferResult = 19
	TransferCodeMustNotBeZero                          CreateTransferResult = 20
	TransferAmountMustNotBeZero                        CreateTransferResult = 21
	TransferDebitAccountNotFound                       CreateTransferResult = 22
	TransferCreditAccountNotFound                      CreateTransferResult = 23
	TransferAccountsMustHaveTheSameLedger              CreateTransferResult = 24
	TransferTransferMustHaveTheSameLedgerAsAccounts    CreateTransferResult = 25
	TransferPendingTransferNotFound                    CreateTransferResult = 26
	TransferPendingTransferNotPending                  CreateTransferResult = 27
	TransferPendingTransferHasDifferentDebitAccountID  CreateTransferResult = 28
	TransferPendingTransferHasDifferentCreditAccountID CreateTransferResult = 29
	TransferPendingTransferHasDifferentLedger          CreateTransferResult = 30
	TransferPendingTransferHasDifferentCode            CreateTransferResult = 31
	TransferExceedsPendingTransferAmount               CreateTransferResult = 32
	TransferPendingTransferHasDifferentAmount          CreateTransferResult = 33
	TransferPendingTransferAlreadyPosted               CreateTransferResult = 34
	TransferPendingTransferAlreadyVoided               CreateTransferResult = 35
	TransferPendingTransferExpired                     CreateTransferResult = 36
	TransferExistsWithDifferentFlags                   CreateTransferResult = 37
	TransferExistsWithDifferentDebitAccountID          CreateTransferResult = 38
	TransferExistsWithDifferentCreditAccountID         CreateTransferResult = 39
	TransferExistsWithDifferentPendingID               CreateTransferResult = 40
	TransferExistsWithDifferentUserData                CreateTransferResult = 41
	TransferExistsWithDifferentTimeout                 CreateTransferResult = 42
	TransferExistsWithDifferentCode                    CreateTransferResult = 43
	TransferExistsWithDifferentAmount                  CreateTransferResult = 44
	TransferExists                                     CreateTransferResult = 45
	TransferOverflowsDebitsPending                     CreateTransferResult = 46
	TransferOverflowsCreditsPending                    CreateTransferResult = 47
	TransferOverflowsDebitsPosted                      CreateTransferResult = 48
	TransferOverflowsCreditsPosted                     CreateTransferResult = 49
	TransferOverflowsDebits                            CreateTransferResult = 50
	TransferOverflowsCredits                           CreateTransferResult = 51
	TransferOverflowsTimeout                           CreateTransferResult = 52
	TransferExceedsCredits                             CreateTransferResult = 53
	TransferExceedsDebits                              CreateTransferResult = 54
)
