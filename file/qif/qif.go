package qif

import "io"
import "os"
import "bufio"
import "time"
import "strings"

import "github.com/peaberberian/GoBanks/database/types"
import "github.com/peaberberian/GoBanks/utils"

func ParseFile(f *os.File, accountDbId int, dateFormat string) (ts []types.Transaction, err error) {
	var label string
	var transactionDate time.Time
	var recordDate time.Time
	var debit float32
	var credit float32
	var reference string

	reader := bufio.NewReader(f)
	for {
		var byteLine []byte
		byteLine, err = reader.ReadBytes('\n')

		if err == io.EOF {
			return ts, nil
		}

		if err != nil {
			return
		}

		if len(byteLine) > 1 {
			strLine := string(byteLine[:len(byteLine)-1])

			switch string(strLine[0]) {

			case "^":
				ts = append(ts, types.Transaction{
					LinkedAccountDbId: accountDbId,
					Label:             label,
					TransactionDate:   transactionDate,
					Debit:             debit,
					Credit:            credit,
					Reference:         reference,
					RecordDate:        recordDate,
				})

				label = ""
				transactionDate = time.Time{}
				recordDate = time.Time{}
				debit = 0
				credit = 0
				reference = ""

			case "!":

			case "D":
				if len(strLine) > 1 {
					dateStr := strLine[1:len(strLine)]
					transactionDate, err = qifStrToDate(dateStr, dateFormat)
					if err != nil {
						return
					}
					recordDate = transactionDate
				}

			case "T":
				if len(strLine) > 1 {
					switch string(strLine[1]) {
					case "-":
						debit, err = utils.StrToFlt32(strLine[2:len(strLine)])
					case "+":
						credit, err = utils.StrToFlt32(strLine[2:len(strLine)])
					default:
						credit, err = utils.StrToFlt32(strLine[1:len(strLine)])
					}
					if err != nil {
						return
					}
				}

			case "N":
				if len(strLine) > 1 {
					reference = strLine[1:len(strLine)]
				}

			case "M":
				if len(strLine) > 1 {
					label = strLine[1:len(strLine)]
				}
			}
		}
	}
}

func qifStrToDate(date string, dateFormat string) (t time.Time, err error) {
	var day string
	var month string
	var year string
	switch dateFormat {
	case "DD/MM/YY":
		strs := strings.Split(date, "/")
		if len(strs) >= 3 {
			day = strs[0]
			month = strs[1]
			year = strs[2]
		}
	case "MM/DD/YY":
		strs := strings.Split(date, "/")
		if len(strs) >= 3 {
			day = strs[1]
			month = strs[0]
			year = strs[2]
		}

	case "YY/MM/DD":
		strs := strings.Split(date, "/")
		if len(strs) >= 3 {
			day = strs[2]
			month = strs[1]
			year = strs[0]
		}
	default:
	}
	if len(year) < 4 {
		year = "20" + year
	}
	t, err = time.Parse("02/01/2006", day+"/"+month+"/"+year)
	if err != nil {
		return
	}
	return
}
