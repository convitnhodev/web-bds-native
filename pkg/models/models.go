package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/deeincom/deeincom/pkg/form"
	"github.com/deeincom/deeincom/pkg/helper"
)

type User struct {
	ID                  int
	Email               string
	Phone               string
	Password            string
	FirstName           string
	LastName            string
	Roles               []string
	EmailToken          string
	PhoneToken          string
	PartnerStatus       string
	LastKYCStatus       string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	SendVerifiedEmailAt time.Time
	SendVerifiedPhoneAt time.Time
	ResetPasswordToken  string
	RPTExpiredAt        *time.Time
}

type Attachment struct {
	ID             int
	Title          string
	ContentType    string
	MineType       string
	Link           string
	Width          int
	Height         int
	Size           int
	VideoLength    int
	VideoThumbnail string
	Product        Product
}

type Product struct {
	ID                   int
	Title                string
	Short                string
	Full                 string
	FullContent          string
	City                 string
	District             string
	Ward                 string
	AddressNumber        string
	Street               string
	HouseDirection       string
	BalconyDirection     string
	BusinessAdvantage    string
	FinancialPlan        string
	Legal                string
	Furniture            string
	Slug                 string
	Type                 string
	PosterLink           string
	HouseCertificateLink string
	FinancePlanLink      string
	Area                 int
	Bedroom              int
	Toilet               int
	Floor                int
	FrontWidth           int
	StreetWidth          int
	PavementWidth        int
	RemainOfSlot         int
	NumOfSlot            int
	CostPerSlot          int
	DepositPercent       float64
	CreatedBy            int
	IsCensorship         bool
	IsSelling            bool
	CensoredAt           *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type Post struct {
	ID          int
	Title       string
	PostType    string
	Slug        string
	Thumbnail   string
	Poster      User
	Tags        []string
	Short       string
	Content     string
	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Comment struct {
	ID            int
	UserId        int
	ParrentId     *int
	Poster        User
	ChildComments []*Comment
	Slug          string
	Message       string
	IsCensorship  bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type File struct {
	ID        int
	LocalPath string
	CloudLink string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type KYC struct {
	ID                int
	UserId            int
	FrontIdentityCard string
	BackIdentityCard  string
	SelfieImage       string
	Feedback          string
	Status            string
	LastKYCFeedback   string
	RejectedBy        *int
	ApprovedBy        *int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Partner struct {
	ID         int
	UserId     int
	Message    string
	CVLink     string
	Status     string
	Feedback   string
	RejectedBy *int
	ApprovedBy *int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Log struct {
	ID        int64
	UserId    int
	Actor     User
	Content   string
	CreatedAt time.Time
}

type Invoice struct {
	ID              int
	UserId          int
	User            User
	Status          string
	InvoiceSerect   string
	TotalAmount     int
	InvoiceSyncedAt *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type InvoiceItem struct {
	ID          int
	InvoiceId   int
	ProductId   int
	Product     Product
	Quatity     int
	CostPerSlot int
	Amount      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Payment struct {
	ID                   int
	InvoiceId            int
	Amount               int
	ActuallyAmount       int
	Status               string
	Method               string
	PayType              string
	TxType               string
	AppotapayTransId     string
	RefundId             string
	RefundResponse       string
	AppotapayAccountNo   string
	AppotapayAccountName string
	AppotapayBankCode    string
	AppotapayBankName    string
	AppotapayBankBranch  string
	TransactionAt        *time.Time
	RefundAt             *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (o *Payment) OrderId(APTPaymentHost string) string {
	if strings.Contains(APTPaymentHost, ".dev.") {
		serectCode := helper.RandString(6)
		return fmt.Sprintf("%s-%d", serectCode, o.ID)
	}
	return fmt.Sprint(o.ID)
}

func (o *Product) Form() *form.Form {
	f := form.New(nil)
	f.Set("Title", o.Title)
	f.Set("Short", o.Short)
	f.Set("Full", o.Full)
	f.Set("FullContent", o.FullContent)
	f.Set("City", o.City)
	f.Set("District", o.District)
	f.Set("Ward", o.Ward)
	f.Set("AddressNumber", o.AddressNumber)
	f.Set("Street", o.Street)
	f.Set("HouseDirection", o.HouseDirection)
	f.Set("BalconyDirection", o.BalconyDirection)
	f.Set("BusinessAdvantage", o.BusinessAdvantage)
	f.Set("FinancialPlan", o.FinancialPlan)
	f.Set("Furniture", o.Furniture)
	f.Set("Type", o.Type)
	f.Set("Legal", o.Legal)

	f.Set("Area", fmt.Sprint(o.Area))
	f.Set("Bedroom", fmt.Sprint(o.Bedroom))
	f.Set("Toilet", fmt.Sprint(o.Toilet))
	f.Set("Floor", fmt.Sprint(o.Floor))
	f.Set("FrontWidth", fmt.Sprint(o.FrontWidth))
	f.Set("StreetWidth", fmt.Sprint(o.StreetWidth))
	f.Set("PavementWidth", fmt.Sprint(o.PavementWidth))

	f.Set("CostPerSlot", fmt.Sprint(o.CostPerSlot))
	f.Set("NumOfSlot", fmt.Sprint(o.NumOfSlot))
	f.Set("DepositPercent", fmt.Sprint(o.DepositPercent))

	return f
}

func (o *Attachment) Form() *form.Form {
	f := form.New(nil)
	f.Set("Title", o.Title)
	f.Set("ContentType", o.ContentType)
	f.Set("MineType", o.MineType)
	f.Set("Link", o.Link)
	f.Set("VideoThumbnail", o.VideoThumbnail)

	f.Set("ProductID", fmt.Sprint(o.Product.ID))
	f.Set("Width", fmt.Sprint(o.Width))
	f.Set("Height", fmt.Sprint(o.Height))
	f.Set("Size", fmt.Sprint(o.Size))
	f.Set("VideoLength", fmt.Sprint(o.VideoLength))

	return f
}

func (o *Post) Form() *form.Form {
	f := form.New(nil)

	f.Set("ID", fmt.Sprint(o.ID))
	f.Set("Title", o.Title)
	f.Set("Content", o.Content)
	f.Set("Tags", strings.Join(o.Tags, ", "))
	f.Set("Short", o.Short)
	f.Set("Thumbnail", o.Thumbnail)
	f.Set("PostType", o.PostType)

	if o.PublishedAt != nil {
		loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
		if err == nil {
			f.Set("PublishedAt", o.PublishedAt.In(loc).Format("2006-01-02T15:04:05"))
		}
	} else {
		f.Set("PublishedAt", "")
	}

	return f
}

func (o *Comment) Form() *form.Form {
	f := form.New(nil)

	f.Set("Comment", o.Message)
	f.Set("Slug", o.Slug)

	return f
}

func (o *KYC) Form() *form.Form {
	f := form.New(nil)

	f.Set("Feedback", o.Feedback)

	return f
}

func (o *Partner) Form() *form.Form {
	f := form.New(nil)

	f.Set("Feedback", o.Feedback)

	return f
}
