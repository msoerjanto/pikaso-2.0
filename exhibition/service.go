package exhibition

import (
	"errors"
	"strconv"
	"strings"

	"github.com/msoerjanto/thepikaso/piece"
	"github.com/msoerjanto/thepikaso/space"
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

type Artist struct {
	ArtistId  int    `json:"artistId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Space struct {
	Location    string `json:"location"`
	SpaceNumber int    `json:"spaceNumber"`
	PieceId     string `json:"pieceId"`
}

type ExhibitionService interface {
	GetPieces() ([]Piece, error)
	CreatePiece(Piece) (Piece, error)
	GetArtists() ([]Artist, error)
	CreateArtist(Artist) (Artist, error)
	GetSpaces() ([]Space, error)
	CreateSpace(Space) (Space, error)
}

type exhibitionService struct {
	pieces  piece.Repository
	artists piece.ArtistRepository
	spaces  space.SpaceRepository
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

func (s *exhibitionService) CreateArtist(new Artist) (Artist, error) {
	toCreate := &piece.Artist{
		ArtistId:  strconv.Itoa(new.ArtistId),
		FirstName: new.FirstName,
		LastName:  new.LastName,
	}

	if err := s.artists.Store(toCreate); err != nil {
		return Artist{}, err
	}

	return new, nil
}

func (s *exhibitionService) CreateSpace(new Space) (Space, error) {
	toCreate := &space.Space{
		SpaceNumber: new.SpaceNumber,
		Location:    new.Location,
	} // omit space id

	if err := s.spaces.Store(toCreate); err != nil {
		return Space{}, err
	}

	return new, nil
}

func (s *exhibitionService) GetPieces() ([]Piece, error) {

	pieces, err := s.pieces.FindAll()
	if err != nil {
		return make([]Piece, 0), err
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

func (s *exhibitionService) GetArtists() ([]Artist, error) {
	artists, err := s.artists.FindAll()
	if err != nil {
		return make([]Artist, 0), err
	}
	var result []Artist
	for _, artist := range artists {
		artistId, err := strconv.Atoi(artist.ArtistId)
		if err != nil {
			return make([]Artist, 0), err
		}
		temp := Artist{
			ArtistId:  artistId,
			FirstName: artist.FirstName,
			LastName:  artist.LastName,
		}
		result = append(result, temp)
	}
	return result, nil
}

func (s *exhibitionService) GetSpaces() ([]Space, error) {
	spaces, err := s.spaces.FindAll()
	if err != nil {
		return make([]Space, 0), err
	}
	var result []Space
	for _, space := range spaces {
		temp := Space{
			Location:    space.Location,
			SpaceNumber: space.SpaceNumber,
			PieceId:     space.PieceId,
		}
		result = append(result, temp)
	}
	return result, nil
}

func NewService(
	pieces piece.Repository,
	artists piece.ArtistRepository,
	spaces space.SpaceRepository) ExhibitionService {
	return &exhibitionService{
		pieces:  pieces,
		artists: artists,
		spaces:  spaces,
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
