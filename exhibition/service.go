package exhibition

import (
	"errors"
	"strconv"
	"strings"

	"github.com/msoerjanto/thepikaso/piece"
)

type Piece struct {
	ArtistId      int    `json:"artistId"`
	PictureNumber int    `json:"pictureNumber"`
	Year          int    `json:"year"`
	Title         string `json:"title"`
	Media         string `json:"media"`
	Length        int    `json:"length"`
	Height        int    `json:"height"`
	Page          int    `json:"page"`
	ImageUrl      string `json:"imageUrl"`
}

type ExhibitionService interface {
	GetPieces() ([]Piece, error)
	CreatePiece(Piece) (Piece, error)
}

type exhibitionService struct {
	pieces piece.Repository
}

func (s *exhibitionService) GetPieces() ([]Piece, error) {

	pieces, err := s.pieces.FindAll()
	if err != nil {
		return make([]Piece, 0), nil
	}
	var result []Piece
	for _, piece := range pieces {
		artistId, picNum, err := getArtistIdAndPicNumFromPieceId(piece.PieceId)
		if err != nil {
			return make([]Piece, 0), nil
		}
		temp := Piece{
			ArtistId:      artistId,
			PictureNumber: picNum,
			Year:          piece.Year,
			Title:         piece.Title,
			Media:         piece.Media,
			Length:        piece.Length,
			Height:        piece.Height,
			Page:          piece.Page,
			ImageUrl:      piece.ImageUrl,
		}
		result = append(result, temp)
	}
	return result, nil
}

func (s *exhibitionService) CreatePiece(new Piece) (Piece, error) {

	toCreate := &piece.Piece{
		PieceId:  strconv.Itoa(new.ArtistId) + "+" + strconv.Itoa(new.PictureNumber),
		Year:     new.Year,
		Title:    new.Title,
		Media:    new.Media,
		Length:   new.Length,
		Height:   new.Height,
		Page:     new.Page,
		ImageUrl: new.ImageUrl,
	}

	if err := s.pieces.Store(toCreate); err != nil {
		return Piece{}, err
	}

	return new, nil
}

func NewService(pieces piece.Repository) ExhibitionService {
	return &exhibitionService{
		pieces: pieces,
	}
}

func getArtistIdAndPicNumFromPieceId(pieceId string) (int, int, error) {
	artistIdAndPicNum := strings.Split(pieceId, "+")
	if len(artistIdAndPicNum) != 2 {
		return 0, 0, errors.New("Bad PieceId")
	}

	artistId, aiderr := strconv.Atoi(artistIdAndPicNum[0])
	picNum, pnumerr := strconv.Atoi(artistIdAndPicNum[1])

	if aiderr != nil {
		return 0, 0, errors.New("Unable to parse artistId")
	} else if pnumerr != nil {
		return 0, 0, errors.New("Unable to parse PictureNumber")
	}

	return artistId, picNum, nil
}
