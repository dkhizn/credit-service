package controller

import (
	"fmt"

	"github.com/shopspring/decimal"
)

func roundToInt(d decimal.Decimal) int {
	i := d.Round(0).IntPart()
	return int(i)
}

func validate(dto ExecuteRequestDTO) error {
	var chosenProgram int
	for _, chosen := range dto.Program {
		if chosen {
			chosenProgram++
		}
	}

	if chosenProgram == 0 {
		return fmt.Errorf("choose program")
	}

	if chosenProgram > 1 {
		return fmt.Errorf("choose only 1 program")
	}

	if dto.InitialPayment*5 < dto.ObjectCost {
		return fmt.Errorf("the initial payment should be more")
	}

	return nil
}
