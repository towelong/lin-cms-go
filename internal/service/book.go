package service

import (
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/towelong/lin-cms-go/internal/domain/dto"
	"github.com/towelong/lin-cms-go/internal/domain/model"
	"github.com/towelong/lin-cms-go/internal/domain/vo"
	"github.com/towelong/lin-cms-go/pkg/response"
	"gorm.io/gorm"
)

type IBookService interface {
	GetBookList() (vos []vo.BookVo)
	FindBookById(id int) (v vo.BookVo, err error)
	DeleteBookById(id int) (err error)
	ChangeBookById(id int, bookDto dto.CreateOrUpdateBookDTO) (err error)
	CreateBook(bookDto dto.CreateOrUpdateBookDTO) (err error)
}

type BookService struct {
	DB *gorm.DB
}

func (b BookService) GetBookList() (vos []vo.BookVo) {
	var books []model.Book
	err := b.DB.Find(&books).Error
	if err != nil {
		vos = make([]vo.BookVo, 0)
	} else {
		copier.Copy(&vos, &books)
	}
	return vos
}

func (b BookService) FindBookById(id int) (v vo.BookVo, err error) {
	var book model.Book
	err = b.DB.First(&book, id).Error
	if err != nil {
		return v, response.New(10022, http.StatusNotFound)
	}
	copier.Copy(&v, &book)
	return v, nil
}

func (b BookService) DeleteBookById(id int) (err error) {
	var book vo.BookVo
	book, err = b.FindBookById(id)
	if err != nil {
		return err
	}
	b.DB.Delete(&model.Book{}, book.ID)
	return nil
}
func (b BookService) ChangeBookById(id int, bookDto dto.CreateOrUpdateBookDTO) (err error) {
	var book vo.BookVo
	book, err = b.FindBookById(id)
	if err != nil {
		return err
	}
	bookModel := model.Book{
		BaseModel: model.BaseModel{ID: book.ID},
	}
	var newBook model.Book
	copier.CopyWithOption(&newBook, &bookDto, copier.Option{IgnoreEmpty: true})
	b.DB.Model(&bookModel).Updates(newBook)
	return nil
}

func (b BookService) CreateBook(bookDto dto.CreateOrUpdateBookDTO) (err error) {
	var book model.Book
	copier.CopyWithOption(&book, &bookDto, copier.Option{IgnoreEmpty: true})
	return b.DB.Omit("create_time, update_time").Create(&book).Error
}
