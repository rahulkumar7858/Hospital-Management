package handlers

import (
	"hospital-management/config"
	"hospital-management/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreatePatient(c *gin.Context) {
	var patient models.Patient
	if err := c.ShouldBindJSON(&patient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	err := config.DB.QueryRow(`
        INSERT INTO patients (name, age, gender, diagnosis) 
        VALUES ($1, $2, $3, $4) 
        RETURNING id
    `, patient.Name, patient.Age, patient.Gender, patient.Diagnosis).Scan(&patient.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create patient"})
		return
	}

	c.JSON(http.StatusOK, patient)
}

func GetPatients(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, name, age, gender, diagnosis FROM patients")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch patients"})
		return
	}
	defer rows.Close()

	patients := []models.Patient{}
	for rows.Next() {
		var p models.Patient
		if err := rows.Scan(&p.ID, &p.Name, &p.Age, &p.Gender, &p.Diagnosis); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Data error"})
			return
		}
		patients = append(patients, p)
	}

	c.JSON(http.StatusOK, patients)
}

func GetPatientByID(c *gin.Context) {
	id := c.Param("id")
	var patient models.Patient
	err := config.DB.QueryRow(`
        SELECT id, name, age, gender, diagnosis 
        FROM patients WHERE id = $1
    `, id).Scan(&patient.ID, &patient.Name, &patient.Age, &patient.Gender, &patient.Diagnosis)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	c.JSON(http.StatusOK, patient)
}

func UpdatePatient(c *gin.Context) {
	id := c.Param("id")
	var patient models.Patient
	if err := c.ShouldBindJSON(&patient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	_, err := config.DB.Exec(`
        UPDATE patients 
        SET name=$1, age=$2, gender=$3, diagnosis=$4 
        WHERE id=$5
    `, patient.Name, patient.Age, patient.Gender, patient.Diagnosis, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update patient"})
		return
	}

	patient.ID, _ = strconv.Atoi(id)
	c.JSON(http.StatusOK, patient)
}

func UpdatePatientDiagnosis(c *gin.Context) {
	id := c.Param("id")
	var update struct {
		Diagnosis string `json:"diagnosis"`
	}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	_, err := config.DB.Exec(`
        UPDATE patients 
        SET diagnosis=$1 
        WHERE id=$2
    `, update.Diagnosis, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update diagnosis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Diagnosis updated successfully"})
}

func DeletePatient(c *gin.Context) {
	id := c.Param("id")
	_, err := config.DB.Exec("DELETE FROM patients WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete patient"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Patient deleted successfully"})
}
