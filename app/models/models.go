package models

import "fmt"

func Models() {
	fmt.Println("Ini models")
}

/*
// endpoint
createPaymentStatus() {
	// create payment 
	https://www.figma.com/proto/nj55wH5f7EkUB3M7kgCgrv/Broker-Insurance-Apps?node-id=175-585&starting-point-node-id=171%3A3
}

// endpoint 
getPaymentStatus() {
	resultPD, _ := db.getAllPaymentData(); // expected: status = pending
	
	// adjustment
	rowAdj, totalAdj := adjustmentCalculation()
	total := resultPD.premiumInception - totalAdj
	
	// payment
	rowPym, balancePym := paymentCalculation(total)


	// paid
	if balancePym >= 0 {
		if resultPD.status != "paid" {
			status := "paid"
			db.save(status)
		} 
	} else { // outstanding
		if resultPD.status == "pending" {
			status := "outstanding"
			db.save(status)
		}
	}

}

// endpoint
createPayment() {
	db.save(pymData) // set data to database
}
// endpoint
createAdjustment() {
	db.save(adjData) // set data to database
}

// function
adjustmentCalculation() {
	resultAdj, _ := db.getAllAdjustment()
	row := resultAdj.length() 
	total := 0
	if (row > 0) {
		for x in resultAdj {
			total = total + x
		}
	}
	return row, total;
}

// function
paymentCalculation(total int) {
	resultPym, _ := db.getAllPayment()
	row := resultPym.length() 
	balance := 0
	if (row > 0) {
		for x in resultPym {
			balance = balance + x
		}
	}
	balance = total - balance 
	return row, balance;
}
*/