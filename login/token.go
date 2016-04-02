package login

import "fmt"
import "time"

import def "github.com/peaberberian/GoBanks/database/definitions"

import jwt "github.com/dgrijalva/jwt-go"

// Representation of the principal attributes of our json web token
type UserToken struct {
	ExpirationDate  time.Time
	UserId          int
	BankIds         []int
	AccountIds      []int
	IsAdministrator bool
}

var tokenSigningKey string

// CreateToken creates a new token string for a specific user.
// Returns an error if a problem with the databases was encountered or if
// the signing key is not secure enough.
func CreateToken(username string, db def.GoBanksDataBase,
) (tokenString string, err error) {
	if tokenSigningKey == "" {
		return "", LoginError{ErrorCode: InvalidSigningKeyError}
	}

	user, err := GetUserFromUsername(db, username)
	if err != nil {
		return "", err
	}
	bnks, err := getBanksForUserId(db, user.DbId)
	if err != nil {
		return "", err
	}
	accs, err := getAccountsForBankIds(db, bnks)
	if err != nil {
		return "", err
	}

	jwToken := jwt.New(jwt.SigningMethodHS256)
	jwToken.Claims["exp"] = time.Now().Add(time.Hour).Unix()
	jwToken.Claims["uid"] = user.DbId
	jwToken.Claims["acc"] = accs
	jwToken.Claims["bnk"] = bnks
	jwToken.Claims["adm"] = user.Administrator
	tokenString, err = jwToken.SignedString([]byte(tokenSigningKey))
	if err != nil {
		return "", LoginError{err: err.Error(),
			ErrorCode: TokenSigningError}
	}

	return tokenString, nil
}

// ParseToken takes in argument the token string and returns an easily
// readable UserToken struct which repeats most of the token properties.
// Returns an error if the token is invalid
func ParseToken(tokenString string) (token UserToken, err error) {

	jwToken, err := jwt.Parse(tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSigningKey), nil
		})

	if err != nil || !jwToken.Valid {
		var err = invalidTokenError{}
		return UserToken{}, LoginError{err: err.Error(),
			ErrorCode: InvalidTokenError}
	}

	if _, ok := jwToken.Method.(*jwt.SigningMethodHMAC); !ok {
		var alg = fmt.Sprintf("%s", jwToken.Header["alg"])
		var err = invalidTokenSigningMethodError{alg}
		return UserToken{}, LoginError{err: err.Error(),
			ErrorCode: InvalidTokenSigningMethodError}
	}

	var expirationDate time.Time
	var userId int
	var bankIds []int
	var accountIds []int
	var isAdministrator bool
	var ok bool

	var createUnreadableError = func(field string) (UserToken, error) {
		err = unreadableTokenError{field: field}
		return UserToken{}, LoginError{err: err.Error(),
			ErrorCode: UnreadableTokenError}
	}

	if userId64, ok := jwToken.Claims["uid"].(float64); !ok {
		return createUnreadableError("uid")
	} else {
		userId = int(userId64)
	}

	if bnkArray, ok := jwToken.Claims["bnk"].([]interface{}); !ok {
		return createUnreadableError("bnk")
	} else {
		for _, val := range bnkArray {
			if bnkId64, ok := val.(float64); !ok {
				return createUnreadableError("bnk")
			} else {
				bankIds = append(bankIds, int(bnkId64))
			}
		}
	}

	if accArray, ok := jwToken.Claims["acc"].([]interface{}); !ok {
		return createUnreadableError("acc")
	} else {
		for _, val := range accArray {
			if accId64, ok := val.(float64); !ok {
				return createUnreadableError("acc")
			} else {
				accountIds = append(accountIds, int(accId64))
			}
		}
	}

	if isAdministrator, ok = jwToken.Claims["adm"].(bool); !ok {
		return createUnreadableError("adm")
	}

	var exp64 float64
	if exp64, ok = jwToken.Claims["exp"].(float64); !ok {
		return createUnreadableError("exp")
	}
	expirationDate = time.Unix(int64(exp64), 0)

	token = UserToken{
		UserId:          userId,
		BankIds:         bankIds,
		AccountIds:      accountIds,
		IsAdministrator: isAdministrator,
		ExpirationDate:  expirationDate,
	}
	if err = verifyToken(token); err != nil {
		return UserToken{}, err
	}

	return token, nil
}

// verifyToken verifies the atm only the expiration date.
// It returns a defined error if it does not pass the checks.
func verifyToken(token UserToken) error {
	if token.ExpirationDate.Unix() <= time.Now().Unix() {
		var err = expiredTokenError{}
		return LoginError{err: err.Error(), ErrorCode: ExpiredTokenError}
	}
	return nil
}

func getAccountsForBankIds(db def.GoBanksDataBase,
	bankIds []int) (accIds []int, err error) {

	var accountsFilter def.BankAccountFilters
	accountsFilter.Filters.Banks = true
	accountsFilter.Values.Banks = bankIds
	accs, err := db.GetBankAccounts(accountsFilter)
	if err != nil {
		return
	}
	for _, acc := range accs {
		accIds = append(accIds, acc.DbId)
	}
	return
}

func getBanksForUserId(db def.GoBanksDataBase,
	userid int) (bnkids []int, err error) {

	var banksFilter def.BankFilters
	banksFilter.Filters.Users = true
	banksFilter.Values.Users = []int{userid}
	bnks, err := db.GetBanks(banksFilter)
	if err != nil {
		return
	}
	for _, bnk := range bnks {
		bnkids = append(bnkids, bnk.DbId)
	}
	return
	return
}

func generateSigningKey() error {
	var err error
	tokenSigningKey, err = generateRandomKey(32)
	return err
}
