package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// Helper: send JSON response
func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// ---------------- Employee Routes ----------------

// GET /employee/{id}
func getEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var emp Emp
	err := employeeCol.FindOne(context.TODO(), bson.M{"id": id}).Decode(&emp)
	if err != nil {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, emp)
}

// POST /employee/{id}/apply-leave
func applyLeaveHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var emp Emp
	if err := employeeCol.FindOne(context.TODO(), bson.M{"id": id}).Decode(&emp); err != nil {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	var l Leave
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	l.Name = emp.Name
	l.EmpID = emp.ID
	l.Status = "Pending"

	_, err := leaveCol.InsertOne(context.TODO(), l)
	if err != nil {
		http.Error(w, "Failed to apply leave", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]string{
		"message": "Leave applied successfully",
		"status":  "Pending",
	})
}

// GET /employee/{id}/leaves
func getEmployeeLeavesHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	cursor, err := leaveCol.Find(context.TODO(), bson.M{"emp_id": id})
	if err != nil {
		http.Error(w, "Error fetching leaves", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var leaves []Leave
	cursor.All(context.TODO(), &leaves)
	jsonResponse(w, leaves)
}

// ---------------- HR Routes ----------------

// GET /hr/pending-leaves
func hrPendingLeavesHandler(w http.ResponseWriter, r *http.Request) {
	cursor, err := leaveCol.Find(context.TODO(), bson.M{"status": "Pending"})
	if err != nil {
		http.Error(w, "Error fetching pending leaves", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var pending []Leave
	cursor.All(context.TODO(), &pending)
	jsonResponse(w, pending)
}

// POST /hr/leave/{index}/grant
func hrGrantLeaveHandler(w http.ResponseWriter, r *http.Request) {
	index := mux.Vars(r)["index"]
	idx, _ := strconv.Atoi(index)

	cursor, _ := leaveCol.Find(context.TODO(), bson.M{"status": "Pending"})
	var leaves []Leave
	cursor.All(context.TODO(), &leaves)

	if idx < 0 || idx >= len(leaves) {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	target := leaves[idx]
	_, err := leaveCol.UpdateOne(context.TODO(),
		bson.M{"emp_id": target.EmpID, "from_date": target.FromDate},
		bson.M{"$set": bson.M{"status": "Granted"}})

	if err != nil {
		http.Error(w, "Failed to update leave", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, map[string]string{"message": "Leave granted successfully"})
}

// POST /hr/leave/{index}/decline
func hrDeclineLeaveHandler(w http.ResponseWriter, r *http.Request) {
	index := mux.Vars(r)["index"]
	idx, _ := strconv.Atoi(index)

	cursor, _ := leaveCol.Find(context.TODO(), bson.M{"status": "Pending"})
	var leaves []Leave
	cursor.All(context.TODO(), &leaves)

	if idx < 0 || idx >= len(leaves) {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	target := leaves[idx]
	_, err := leaveCol.UpdateOne(context.TODO(),
		bson.M{"emp_id": target.EmpID, "from_date": target.FromDate},
		bson.M{"$set": bson.M{"status": "Declined"}})

	if err != nil {
		http.Error(w, "Failed to update leave", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, map[string]string{"message": "Leave declined successfully"})
}
