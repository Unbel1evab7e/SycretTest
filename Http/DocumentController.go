package Http

import (
	"SycretTest/Models"
	"SycretTest/Models/Document"
	"SycretTest/Services"
	"SycretTest/httputil"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"path"
)

type DocumentController struct {
	DocService *Services.DocumentService
}

func NewDocumentController(r *gin.Engine, docService *Services.DocumentService) {
	controller := &DocumentController{DocService: docService}

	v1 := r.Group("api/v1")

	v1.POST("/Documents", controller.GetFinalDocument)
}

// GetFinalDocument godoc
// @Summary      GetFinalDocument
// @Description  Return Created documentUrl
// @Tags         Documents
// @Accept       json
// @Produce      json
// @Param        doc  body      Document.DocumentRequest true  "Get Document"
// @Success      201     {object}  string
// @Failure      422     {object}  httputil.HTTPError
// @Failure      400     {object}  httputil.HTTPError
// @Router       /Documents [post]
func (dc DocumentController) GetFinalDocument(c *gin.Context) {
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		httputil.NewError(c, http.StatusUnprocessableEntity, Models.ErrBadParamInput)
		return
	}

	var body Document.DocumentRequest

	err = json.Unmarshal(b, &body)
	if err != nil {
		httputil.NewError(c, http.StatusUnprocessableEntity, Models.ErrBadParamInput)
		return
	}
	err, valRes := validateRequest(&body)
	if err != nil || !valRes {
		httputil.NewError(c, http.StatusBadRequest, Models.ErrBadParamInput)
		return
	}

	docx, err := dc.DocService.GetDocument(&body)
	if err != nil || !valRes {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}
	url := c.Request.Host
	response := Document.DocumentResponse{
		URLWord: path.Join(url, docx),
	}

	c.JSON(http.StatusOK, response)

}

func validateRequest(request *Document.DocumentRequest) (error, bool) {
	_, err := url.Parse(request.URLTemplate)

	if err != nil {
		return err, false
	}

	if request.RecordID == 0 {
		return errors.New("record Id is nil"), false
	}

	return nil, true
}
