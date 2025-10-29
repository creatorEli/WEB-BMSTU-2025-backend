package handler

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	_ "time_of_armies/docs"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

// MinIO клиент
var minioClient *minio.Client

// Конфигурация
const (
	minioEndpoint  = "localhost:9000"
	minioAccessKey = "minio"    // замените на свои
	minioSecretKey = "minio124" // замените на свои
	minioUseSSL    = false
	bucketName     = "armies"
)

func InitMinIO() error {
	var err error
	minioClient, err = minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKey, minioSecretKey, ""),
		Secure: minioUseSSL,
	})
	if err != nil {
		return err
	}
	// Проверяем существование бакета
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	logrus.Info(exists)

	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		log.Printf("Bucket %s created successfully", bucketName)
	}

	return nil
}

type UploadResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	FileName string `json:"file_name"`
	FileURL  string `json:"file_url"`
	FileSize int64  `json:"file_size"`
}

// Структура для загрузки
type ImageUploadRequest struct {
	Image *multipart.FileHeader `form:"image_army" binding:"required"`
}

// uploadArmyImage godoc
// @Summary		 Добавить или Обновить изображение армии
// @Description  Добавить или Обновить изображение армии зная её идентификатор
// @Tags         Requests
// @Produce      json
// Param [name] [type] [dataType] [required] [description]
// @Param		 id formData int true "ID армии к которой добавляется изображение"
// @Param		 image_army formData file true "Изображение армии"
// @Success      200  {object} UploadResponse
// @Router       /army/{id}/upload_image [post]
func (h *Handler) uploadArmyImage(c *gin.Context) {
	logrus.Info("we are here UAI 1!")

	var req ImageUploadRequest

	// Валидация формы
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, UploadResponse{
			Success: false,
			Message: "Неверный запрос: " + err.Error(),
		})
		return
	}

	// Проверка типа файла
	if !isImageFile(req.Image) {
		c.JSON(http.StatusBadRequest, UploadResponse{
			Success: false,
			Message: "Файл должен быть картинкой (JPEG, PNG, GIF)",
		})
		return
	}
	// Проверка размера файла (максимум 10MB)
	if req.Image.Size > 10<<20 {
		c.JSON(http.StatusBadRequest, UploadResponse{
			Success: false,
			Message: "Слишком большой размер файла. Максимальный размер 10 МБ",
		})
		return
	}

	fileURL, err := uploadToMinIO(c, req.Image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UploadResponse{
			Success: false,
			Message: "Не удалось загрузить изображение: " + err.Error(),
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	// НАДО В УСЛУГУ ПО ЕЕ ID ДОБАВЛЯТЬ КАРТИНКУ

	err = h.Repository.SetArmyImage(id, fileURL)
	logrus.Info("we are here 10!")
	if err != nil {
		c.JSON(http.StatusInternalServerError, UploadResponse{
			Success: false,
			Message: "Не удалось установить URI изображения в базе данных: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, UploadResponse{
		Success:  true,
		Message:  "Изображение успешно добавлено",
		FileName: req.Image.Filename,
		FileURL:  fileURL,
		FileSize: req.Image.Size,
	})
}

// uploadToMinIO загружает файл в MinIO
func uploadToMinIO(c *gin.Context, fileHeader *multipart.FileHeader) (string, error) {
	// Открываем файл
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()
	logrus.Info("we got here 1!")
	// Генерируем уникальное имя файла
	fileExt := filepath.Ext(fileHeader.Filename)
	fileName := generateFileName(fileExt)
	logrus.Info("we got here 2!" + fileName)

	// Загружаем в MinIO
	ctx := context.Background()
	logrus.Info("we got here 3!")

	err = InitMinIO()
	if err != nil {
		return "", err
	}

	info, err := minioClient.PutObject(ctx, bucketName, fileName, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: fileHeader.Header.Get("Content-Type"),
	})
	logrus.Info("we got here 4!")

	if err != nil {
		return "", err
	}

	log.Printf("Successfully uploaded %s of size %d\n", fileName, info.Size)

	// Генерируем URL для доступа к файлу
	fileURL := fmt.Sprintf("http://%s/%s/%s", minioEndpoint, bucketName, fileName)
	return fileURL, nil
}

// generateFileName создает уникальное имя файла
func generateFileName(extension string) string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("army_%d%s", timestamp, extension)
}

// isImageFile проверяет, является ли файл изображением
func isImageFile(fileHeader *multipart.FileHeader) bool {
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	contentType := fileHeader.Header.Get("Content-Type")
	return allowedTypes[contentType]
}
