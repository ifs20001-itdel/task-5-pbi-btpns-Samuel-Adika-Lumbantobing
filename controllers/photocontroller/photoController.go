package photoController

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/jeypc/go-jwt-mux/helper"

	"github.com/jeypc/go-jwt-mux/models"
)

func CreatePhoto(w http.ResponseWriter, r *http.Request) {
	var photoInput models.Photo
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&photoInput); err != nil {
		response := map[string]string{"message": "Invalid request body"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// Validasi bahwa UserID yang diberikan sudah ada di database
	var user models.User
	if err := models.DB.First(&user, photoInput.UserID).Error; err != nil {
		response := map[string]string{"message": "User not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	// Simpan foto ke database
	if err := models.DB.Create(&photoInput).Error; err != nil {
		response := map[string]string{"message": "Error creating photo"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Photo created successfully"}
	helper.ResponseJSON(w, http.StatusCreated, response)
}

func GetPhotos(w http.ResponseWriter, r *http.Request) {
	var photos []models.Photo

	// Ambil semua foto dari database
	if err := models.DB.Find(&photos).Error; err != nil {
		response := map[string]string{"message": "Error fetching photos"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	helper.ResponseJSON(w, http.StatusOK, photos)
}

func UpdatePhoto(w http.ResponseWriter, r *http.Request) {
	photoID := chi.URLParam(r, "photoID")

	var photo models.Photo

	// Ambil informasi foto dari database berdasarkan photoID
	if err := models.DB.First(&photo, photoID).Error; err != nil {
		response := map[string]string{"message": "Photo not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	// Periksa apakah pengguna yang meminta update adalah pemilik foto
	if !isAuthorized(r, photo.UserID) {
		response := map[string]string{"message": "Unauthorized"}
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	var photoInput models.Photo
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&photoInput); err != nil {
		response := map[string]string{"message": "Invalid request body"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// Update foto di database berdasarkan photoID
	if err := models.DB.Model(&models.Photo{}).Where("id = ?", photoID).Updates(&photoInput).Error; err != nil {
		response := map[string]string{"message": "Error updating photo"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Photo updated successfully"}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func DeletePhoto(w http.ResponseWriter, r *http.Request) {
	photoID := chi.URLParam(r, "photoID")

	var photo models.Photo

	// Ambil informasi foto dari database berdasarkan photoID
	if err := models.DB.First(&photo, photoID).Error; err != nil {
		response := map[string]string{"message": "Photo not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	// Periksa apakah pengguna yang meminta delete adalah pemilik foto
	if !isAuthorized(r, photo.UserID) {
		response := map[string]string{"message": "Unauthorized"}
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	// Hapus foto dari database berdasarkan photoID
	if err := models.DB.Where("id = ?", photoID).Delete(&models.Photo{}).Error; err != nil {
		response := map[string]string{"message": "Error deleting photo"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Photo deleted successfully"}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func UploadPhoto(w http.ResponseWriter, r *http.Request) {
	// Ambil ID user dari path parameter
	userID := chi.URLParam(r, "userID")

	// Ambil file dari form-data
	file, header, err := r.FormFile("photo")
	if err != nil {
		response := map[string]string{"message": "Error getting file from form-data"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer file.Close()

	// Simpan file ke server (misalnya di folder "uploads")
	filename := userID + "_" + header.Filename
	destination := filepath.Join("uploads", filename)
	if err := helper.SaveFile(file, destination); err != nil {
		response := map[string]string{"message": "Error saving file"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// Simpan informasi foto ke database
	photo := models.Photo{
		UserID: userID,
		URL:    destination,
	}

	if err := models.DB.Create(&photo).Error; err != nil {
		response := map[string]string{"message": "Error creating photo"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Photo uploaded successfully"}
	helper.ResponseJSON(w, http.StatusCreated, response)
}

func GetPhotosByUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	var photos []models.Photo

	// Ambil semua foto dari database berdasarkan userID
	if err := models.DB.Where("user_id = ?", userID).Find(&photos).Error; err != nil {
		response := map[string]string{"message": "Error fetching photos"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	helper.ResponseJSON(w, http.StatusOK, photos)
}

// Fungsi isAuthorized untuk memeriksa apakah pengguna yang meminta tindakan adalah pemilik foto
func isAuthorized(r *http.Request, photoUserID string) bool {
	requestedUserID := r.Header.Get("X-UserID")
	return requestedUserID == photoUserID
}
