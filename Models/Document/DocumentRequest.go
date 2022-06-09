package Document

type DocumentRequest struct {
	URLTemplate string `json:"url_template" example:"https://sycret.ru/service/apigendoc/forma_025u.doc"`
	RecordID    int    `json:"record_id" example:"30" format:"int"`
}
