package controllers

import (
	"io"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductImgController struct {
	client *mongo.Client
}

func NewProductImgController(c *mongo.Client) *ProductImgController {
	return &ProductImgController{c}
}

func (pic ProductImgController) UploadProductImg(w http.ResponseWriter , req *http.Request , _ httprouter.Params) {

	if req.Method == http.MethodPost{
		file , handler , err := req.FormFile("productImg") // Form'dan gelen dosyayı alıyoruz.
		if err != nil {
			http.Error(w,err.Error(),http.StatusBadRequest)
			return
		}
		defer file.Close() 

		out , err := os.Create("./media/products/"+handler.Filename) // Kullanıcının yüklediği dosyayı media/products klasörüne kaydetme. İsim olarak kullanıcıdan alıyoruz filename'i bu sonradan değişebilir. Kendimiz'de verebiliriz.
		if err != nil {
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		defer out.Close()

		_ , err = io.Copy(out,file) // Aldıgımız dosyayı , belirlediğimiz yere kopyalama.
		if err != nil {
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Image uploaded successfully"))


	}else{
		http.Error(w,"Method Not Allowed",http.StatusMethodNotAllowed)
		return
	}
	

}
