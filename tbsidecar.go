package tbsidecar

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"lukechampine.com/uint128"
)

func BytesToHexString(bs [16]byte) string {
	newBs := make([]byte, 32) // twice the size for hex encoding
	hex.Encode(newBs, bs[:])
	return string(newBs)
}

func HexStringToBytes(s string) ([16]byte, error) {
	if len(s) != 32 {
		return [16]byte{}, errors.New("input string is not the correct length for hex decoding")
	}
	newBs := make([]byte, 16)
	_, err := hex.Decode(newBs, []byte(s))
	if err != nil {
		return [16]byte{}, errors.New("failed to decode string into bytes")
	}
	var arr [16]byte
	copy(arr[:], newBs)
	return arr, nil
}

func BytesToString(bs [16]byte) string {
	var bsa []byte
	copy(bsa, bs[:])
	return uint128.FromBytes(bsa).String()
}

type AccountFlags struct {
	Linked                     bool
	DebitsMustNotExceedCredits bool
	CreditsMustNotExceedDebits bool
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

type Accounts struct {
	Accounts []Account
}

type Account struct {
	ID             [16]byte
	DebitsPending  [16]byte
	DebitsPosted   [16]byte
	CreditsPending [16]byte
	CreditsPosted  [16]byte
	UserData128    [16]byte
	UserData64     uint64
	UserData32     uint32
	Reserved       uint32
	Ledger         uint32
	Code           uint16
	Flags          uint16
	Timestamp      uint64
}

func (a Account) String() string {
	return fmt.Sprintf(
		"id: %s\ndebits_pending: %s\ndebits_posted: %s\ncredits_pending: %s\ncredits_posted: %s\nledger: %d\nflags: %d",
		uuid.UUID(a.ID),
		BytesToString(a.DebitsPending),
		BytesToString(a.DebitsPosted),
		BytesToString(a.CreditsPending),
		BytesToString(a.CreditsPosted),
		a.Ledger,
		a.Flags,
	)
}

func (o Account) AccountFlags() AccountFlags {
	var f AccountFlags
	f.Linked = ((o.Flags >> 0) & 0x1) == 1
	f.DebitsMustNotExceedCredits = ((o.Flags >> 1) & 0x1) == 1
	f.CreditsMustNotExceedDebits = ((o.Flags >> 2) & 0x1) == 1
	return f
}

type Transfers struct {
	Transfers []Transfer
}

type Transfer struct {
	ID              [16]byte
	DebitAccountID  [16]byte
	CreditAccountID [16]byte
	Amount          [16]byte
	PendingID       [16]byte
	UserData128     [16]byte
	UserData64      uint64
	UserData32      uint32
	Timeout         uint32
	Ledger          uint32
	Code            uint16
	Flags           uint16
	Timestamp       uint64
}

func (t Transfer) String() string {
	return fmt.Sprintf(
		"id: %s\ndebit_account: %s\ncredit_account: %s\nledger: %d\namount: %s",
		uuid.UUID(t.ID),
		uuid.UUID(t.DebitAccountID),
		uuid.UUID(t.CreditAccountID),
		t.Ledger,
		BytesToString(t.Amount),
	)
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

type CreateAccountResult uint32

const (
	AccountOK                             CreateAccountResult = 0
	AccountLinkedEventFailed              CreateAccountResult = 1
	AccountLinkedEventChainOpen           CreateAccountResult = 2
	AccountTimestampMustBeZero            CreateAccountResult = 3
	AccountReservedField                  CreateAccountResult = 4
	AccountReservedFlag                   CreateAccountResult = 5
	AccountIDMustNotBeZero                CreateAccountResult = 6
	AccountIDMustNotBeIntMax              CreateAccountResult = 7
	AccountFlagsAreMutuallyExclusive      CreateAccountResult = 8
	AccountDebitsPendingMustBeZero        CreateAccountResult = 9
	AccountDebitsPostedMustBeZero         CreateAccountResult = 10
	AccountCreditsPendingMustBeZero       CreateAccountResult = 11
	AccountCreditsPostedMustBeZero        CreateAccountResult = 12
	AccountLedgerMustNotBeZero            CreateAccountResult = 13
	AccountCodeMustNotBeZero              CreateAccountResult = 14
	AccountExistsWithDifferentFlags       CreateAccountResult = 15
	AccountExistsWithDifferentUserData128 CreateAccountResult = 16
	AccountExistsWithDifferentUserData64  CreateAccountResult = 17
	AccountExistsWithDifferentUserData32  CreateAccountResult = 18
	AccountExistsWithDifferentLedger      CreateAccountResult = 19
	AccountExistsWithDifferentCode        CreateAccountResult = 20
	AccountExists                         CreateAccountResult = 21
)

func (i CreateAccountResult) String() string {
	switch i {
	case AccountOK:
		return "AccountOK"
	case AccountLinkedEventFailed:
		return "AccountLinkedEventFailed"
	case AccountLinkedEventChainOpen:
		return "AccountLinkedEventChainOpen"
	case AccountTimestampMustBeZero:
		return "AccountTimestampMustBeZero"
	case AccountReservedField:
		return "AccountReservedField"
	case AccountReservedFlag:
		return "AccountReservedFlag"
	case AccountIDMustNotBeZero:
		return "AccountIDMustNotBeZero"
	case AccountIDMustNotBeIntMax:
		return "AccountIDMustNotBeIntMax"
	case AccountFlagsAreMutuallyExclusive:
		return "AccountFlagsAreMutuallyExclusive"
	case AccountDebitsPendingMustBeZero:
		return "AccountDebitsPendingMustBeZero"
	case AccountDebitsPostedMustBeZero:
		return "AccountDebitsPostedMustBeZero"
	case AccountCreditsPendingMustBeZero:
		return "AccountCreditsPendingMustBeZero"
	case AccountCreditsPostedMustBeZero:
		return "AccountCreditsPostedMustBeZero"
	case AccountLedgerMustNotBeZero:
		return "AccountLedgerMustNotBeZero"
	case AccountCodeMustNotBeZero:
		return "AccountCodeMustNotBeZero"
	case AccountExistsWithDifferentFlags:
		return "AccountExistsWithDifferentFlags"
	case AccountExistsWithDifferentUserData128:
		return "AccountExistsWithDifferentUserData128"
	case AccountExistsWithDifferentUserData64:
		return "AccountExistsWithDifferentUserData64"
	case AccountExistsWithDifferentUserData32:
		return "AccountExistsWithDifferentUserData32"
	case AccountExistsWithDifferentLedger:
		return "AccountExistsWithDifferentLedger"
	case AccountExistsWithDifferentCode:
		return "AccountExistsWithDifferentCode"
	case AccountExists:
		return "AccountExists"
	}
	return "CreateAccountResult(" + strconv.FormatInt(int64(i+1), 10) + ")"
}

type CreateTransferResult uint32

const (
	TransferOK                                         CreateTransferResult = 0
	TransferLinkedEventFailed                          CreateTransferResult = 1
	TransferLinkedEventChainOpen                       CreateTransferResult = 2
	TransferTimestampMustBeZero                        CreateTransferResult = 3
	TransferReservedFlag                               CreateTransferResult = 4
	TransferIDMustNotBeZero                            CreateTransferResult = 5
	TransferIDMustNotBeIntMax                          CreateTransferResult = 6
	TransferFlagsAreMutuallyExclusive                  CreateTransferResult = 7
	TransferDebitAccountIDMustNotBeZero                CreateTransferResult = 8
	TransferDebitAccountIDMustNotBeIntMax              CreateTransferResult = 9
	TransferCreditAccountIDMustNotBeZero               CreateTransferResult = 10
	TransferCreditAccountIDMustNotBeIntMax             CreateTransferResult = 11
	TransferAccountsMustBeDifferent                    CreateTransferResult = 12
	TransferPendingIDMustBeZero                        CreateTransferResult = 13
	TransferPendingIDMustNotBeZero                     CreateTransferResult = 14
	TransferPendingIDMustNotBeIntMax                   CreateTransferResult = 15
	TransferPendingIDMustBeDifferent                   CreateTransferResult = 16
	TransferTimeoutReservedForPendingTransfer          CreateTransferResult = 17
	TransferAmountMustNotBeZero                        CreateTransferResult = 18
	TransferLedgerMustNotBeZero                        CreateTransferResult = 19
	TransferCodeMustNotBeZero                          CreateTransferResult = 20
	TransferDebitAccountNotFound                       CreateTransferResult = 21
	TransferCreditAccountNotFound                      CreateTransferResult = 22
	TransferAccountsMustHaveTheSameLedger              CreateTransferResult = 23
	TransferTransferMustHaveTheSameLedgerAsAccounts    CreateTransferResult = 24
	TransferPendingTransferNotFound                    CreateTransferResult = 25
	TransferPendingTransferNotPending                  CreateTransferResult = 26
	TransferPendingTransferHasDifferentDebitAccountID  CreateTransferResult = 27
	TransferPendingTransferHasDifferentCreditAccountID CreateTransferResult = 28
	TransferPendingTransferHasDifferentLedger          CreateTransferResult = 29
	TransferPendingTransferHasDifferentCode            CreateTransferResult = 30
	TransferExceedsPendingTransferAmount               CreateTransferResult = 31
	TransferPendingTransferHasDifferentAmount          CreateTransferResult = 32
	TransferPendingTransferAlreadyPosted               CreateTransferResult = 33
	TransferPendingTransferAlreadyVoided               CreateTransferResult = 34
	TransferPendingTransferExpired                     CreateTransferResult = 35
	TransferExistsWithDifferentFlags                   CreateTransferResult = 36
	TransferExistsWithDifferentDebitAccountID          CreateTransferResult = 37
	TransferExistsWithDifferentCreditAccountID         CreateTransferResult = 38
	TransferExistsWithDifferentAmount                  CreateTransferResult = 39
	TransferExistsWithDifferentPendingID               CreateTransferResult = 40
	TransferExistsWithDifferentUserData128             CreateTransferResult = 41
	TransferExistsWithDifferentUserData64              CreateTransferResult = 42
	TransferExistsWithDifferentUserData32              CreateTransferResult = 43
	TransferExistsWithDifferentTimeout                 CreateTransferResult = 44
	TransferExistsWithDifferentCode                    CreateTransferResult = 45
	TransferExists                                     CreateTransferResult = 46
	TransferOverflowsDebitsPending                     CreateTransferResult = 47
	TransferOverflowsCreditsPending                    CreateTransferResult = 48
	TransferOverflowsDebitsPosted                      CreateTransferResult = 49
	TransferOverflowsCreditsPosted                     CreateTransferResult = 50
	TransferOverflowsDebits                            CreateTransferResult = 51
	TransferOverflowsCredits                           CreateTransferResult = 52
	TransferOverflowsTimeout                           CreateTransferResult = 53
	TransferExceedsCredits                             CreateTransferResult = 54
	TransferExceedsDebits                              CreateTransferResult = 55
)

func (i CreateTransferResult) String() string {
	switch i {
	case TransferOK:
		return "TransferOK"
	case TransferLinkedEventFailed:
		return "TransferLinkedEventFailed"
	case TransferLinkedEventChainOpen:
		return "TransferLinkedEventChainOpen"
	case TransferTimestampMustBeZero:
		return "TransferTimestampMustBeZero"
	case TransferReservedFlag:
		return "TransferReservedFlag"
	case TransferIDMustNotBeZero:
		return "TransferIDMustNotBeZero"
	case TransferIDMustNotBeIntMax:
		return "TransferIDMustNotBeIntMax"
	case TransferFlagsAreMutuallyExclusive:
		return "TransferFlagsAreMutuallyExclusive"
	case TransferDebitAccountIDMustNotBeZero:
		return "TransferDebitAccountIDMustNotBeZero"
	case TransferDebitAccountIDMustNotBeIntMax:
		return "TransferDebitAccountIDMustNotBeIntMax"
	case TransferCreditAccountIDMustNotBeZero:
		return "TransferCreditAccountIDMustNotBeZero"
	case TransferCreditAccountIDMustNotBeIntMax:
		return "TransferCreditAccountIDMustNotBeIntMax"
	case TransferAccountsMustBeDifferent:
		return "TransferAccountsMustBeDifferent"
	case TransferPendingIDMustBeZero:
		return "TransferPendingIDMustBeZero"
	case TransferPendingIDMustNotBeZero:
		return "TransferPendingIDMustNotBeZero"
	case TransferPendingIDMustNotBeIntMax:
		return "TransferPendingIDMustNotBeIntMax"
	case TransferPendingIDMustBeDifferent:
		return "TransferPendingIDMustBeDifferent"
	case TransferTimeoutReservedForPendingTransfer:
		return "TransferTimeoutReservedForPendingTransfer"
	case TransferAmountMustNotBeZero:
		return "TransferAmountMustNotBeZero"
	case TransferLedgerMustNotBeZero:
		return "TransferLedgerMustNotBeZero"
	case TransferCodeMustNotBeZero:
		return "TransferCodeMustNotBeZero"
	case TransferDebitAccountNotFound:
		return "TransferDebitAccountNotFound"
	case TransferCreditAccountNotFound:
		return "TransferCreditAccountNotFound"
	case TransferAccountsMustHaveTheSameLedger:
		return "TransferAccountsMustHaveTheSameLedger"
	case TransferTransferMustHaveTheSameLedgerAsAccounts:
		return "TransferTransferMustHaveTheSameLedgerAsAccounts"
	case TransferPendingTransferNotFound:
		return "TransferPendingTransferNotFound"
	case TransferPendingTransferNotPending:
		return "TransferPendingTransferNotPending"
	case TransferPendingTransferHasDifferentDebitAccountID:
		return "TransferPendingTransferHasDifferentDebitAccountID"
	case TransferPendingTransferHasDifferentCreditAccountID:
		return "TransferPendingTransferHasDifferentCreditAccountID"
	case TransferPendingTransferHasDifferentLedger:
		return "TransferPendingTransferHasDifferentLedger"
	case TransferPendingTransferHasDifferentCode:
		return "TransferPendingTransferHasDifferentCode"
	case TransferExceedsPendingTransferAmount:
		return "TransferExceedsPendingTransferAmount"
	case TransferPendingTransferHasDifferentAmount:
		return "TransferPendingTransferHasDifferentAmount"
	case TransferPendingTransferAlreadyPosted:
		return "TransferPendingTransferAlreadyPosted"
	case TransferPendingTransferAlreadyVoided:
		return "TransferPendingTransferAlreadyVoided"
	case TransferPendingTransferExpired:
		return "TransferPendingTransferExpired"
	case TransferExistsWithDifferentFlags:
		return "TransferExistsWithDifferentFlags"
	case TransferExistsWithDifferentDebitAccountID:
		return "TransferExistsWithDifferentDebitAccountID"
	case TransferExistsWithDifferentCreditAccountID:
		return "TransferExistsWithDifferentCreditAccountID"
	case TransferExistsWithDifferentAmount:
		return "TransferExistsWithDifferentAmount"
	case TransferExistsWithDifferentPendingID:
		return "TransferExistsWithDifferentPendingID"
	case TransferExistsWithDifferentUserData128:
		return "TransferExistsWithDifferentUserData128"
	case TransferExistsWithDifferentUserData64:
		return "TransferExistsWithDifferentUserData64"
	case TransferExistsWithDifferentUserData32:
		return "TransferExistsWithDifferentUserData32"
	case TransferExistsWithDifferentTimeout:
		return "TransferExistsWithDifferentTimeout"
	case TransferExistsWithDifferentCode:
		return "TransferExistsWithDifferentCode"
	case TransferExists:
		return "TransferExists"
	case TransferOverflowsDebitsPending:
		return "TransferOverflowsDebitsPending"
	case TransferOverflowsCreditsPending:
		return "TransferOverflowsCreditsPending"
	case TransferOverflowsDebitsPosted:
		return "TransferOverflowsDebitsPosted"
	case TransferOverflowsCreditsPosted:
		return "TransferOverflowsCreditsPosted"
	case TransferOverflowsDebits:
		return "TransferOverflowsDebits"
	case TransferOverflowsCredits:
		return "TransferOverflowsCredits"
	case TransferOverflowsTimeout:
		return "TransferOverflowsTimeout"
	case TransferExceedsCredits:
		return "TransferExceedsCredits"
	case TransferExceedsDebits:
		return "TransferExceedsDebits"
	}
	return "CreateTransferResult(" + strconv.FormatInt(int64(i+1), 10) + ")"
}

type AccountEventResult struct {
	Index  uint32
	Result CreateAccountResult
}

type TransferEventResult struct {
	Index  uint32
	Result CreateTransferResult
}
