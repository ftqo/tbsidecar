package tbsidecar

import "encoding/hex"

func BytesToString(bs [16]byte) string {
	newBs := make([]byte, 16)
	hex.Encode(newBs, bs[:])
	return string(newBs)
}

func StringToBytes(s string) [16]byte {
	newBs := make([]byte, 16)
	hex.Decode(newBs, []byte(s))
	return [16]byte(newBs)
}

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
	AccountOK CreateAccountResult = iota
	AccountLinkedEventFailed
	AccountLinkedEventChainOpen
	AccountTimestampMustBeZero
	AccountReservedFlag
	AccountReservedField
	AccountIDMustNotBeZero
	AccountIDMustNotBeIntMax
	AccountFlagsAreMutuallyExclusive
	AccountLedgerMustNotBeZero
	AccountCodeMustNotBeZero
	AccountDebitsPendingMustBeZero
	AccountDebitsPostedMustBeZero
	AccountCreditsPendingMustBeZero
	AccountCreditsPostedMustBeZero
	AccountExistsWithDifferentFlags
	AccountExistsWithDifferentUserData
	AccountExistsWithDifferentLedger
	AccountExistsWithDifferentCode
	AccountExists
)

type CreateTransferResult uint32

const (
	TransferOK CreateTransferResult = iota
	TransferLinkedEventFailed
	TransferLinkedEventChainOpen
	TransferTimestampMustBeZero
	TransferReservedFlag
	TransferReservedField
	TransferIDMustNotBeZero
	TransferIDMustNotBeIntMax
	TransferFlagsAreMutuallyExclusive
	TransferDebitAccountIDMustNotBeZero
	TransferDebitAccountIDMustNotBeIntMax
	TransferCreditAccountIDMustNotBeZero
	TransferCreditAccountIDMustNotBeIntMax
	TransferAccountsMustBeDifferent
	TransferPendingIDMustBeZero
	TransferPendingIDMustNotBeZero
	TransferPendingIDMustNotBeIntMax
	TransferPendingIDMustBeDifferent
	TransferTimeoutReservedForPendingTransfer
	TransferLedgerMustNotBeZero
	TransferCodeMustNotBeZero
	TransferAmountMustNotBeZero
	TransferDebitAccountNotFound
	TransferCreditAccountNotFound
	TransferAccountsMustHaveTheSameLedger
	TransferTransferMustHaveTheSameLedgerAsAccounts
	TransferPendingTransferNotFound
	TransferPendingTransferNotPending
	TransferPendingTransferHasDifferentDebitAccountID
	TransferPendingTransferHasDifferentCreditAccountID
	TransferPendingTransferHasDifferentLedger
	TransferPendingTransferHasDifferentCode
	TransferExceedsPendingTransferAmount
	TransferPendingTransferHasDifferentAmount
	TransferPendingTransferAlreadyPosted
	TransferPendingTransferAlreadyVoided
	TransferPendingTransferExpired
	TransferExistsWithDifferentFlags
	TransferExistsWithDifferentDebitAccountID
	TransferExistsWithDifferentCreditAccountID
	TransferExistsWithDifferentPendingID
	TransferExistsWithDifferentUserData
	TransferExistsWithDifferentTimeout
	TransferExistsWithDifferentCode
	TransferExistsWithDifferentAmount
	TransferExists
	TransferOverflowsDebitsPending
	TransferOverflowsCreditsPending
	TransferOverflowsDebitsPosted
	TransferOverflowsCreditsPosted
	TransferOverflowsDebits
	TransferOverflowsCredits
	TransferOverflowsTimeout
	TransferExceedsCredits
	TransferExceedsDebits
)
