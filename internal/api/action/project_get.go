package action

import (
	"github.com/gbooksoft/backend/internal/domain/project/application"
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
)

type ProjectGet struct {
}

func NewProjectGet() *ProjectGet {
	return new(ProjectGet)
}

type chapterResponse struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Order   int    `json:"order"`
	Version int    `json:"version"`
}

type characterResponse struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Description *string           `json:"description"`
	Attributes  map[string]string `json:"attributes"`
}

type locationResponse struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type metadataResponse struct {
	Author   string `json:"author"`
	Keywords string `json:"keywords"`
	Language string `json:"language"`
	Isbn     string `json:"isbn"`
}
type projectFullItemResponse struct {
	Id          string                     `json:"id"`
	Title       string                     `json:"title"`
	Description string                     `json:"description"`
	Genres      []projectGenreItemResponse `json:"genres"`
	Chapters    []chapterResponse          `json:"chapters"`
	Characters  []characterResponse        `json:"characters"`
	Locations   []locationResponse         `json:"locations"`
	Metadata    metadataResponse           `json:"metadata"`
}

func (a *ProjectGet) Invoke(ctx *gin.Context) {
	claims, err := getJwtClaims(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	projectId := ctx.Params.ByName("id")

	query := &application.GetProjectQuery{
		ProjectId: projectId,
		AuthorId:  claims.Subject,
	}

	resp, err := mediatr.Send[*application.GetProjectQuery, *application.GetProjectHandlerResponse](ctx, query)

	if err != nil {
		ctx.Error(err)
		return
	}

	if resp == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	responseItem := projectFullItemResponse{
		Id:          resp.Project.Id,
		Title:       resp.Project.Title,
		Description: resp.Project.Description,
		Metadata: metadataResponse{
			Author:   resp.Project.Metadata.Author,
			Keywords: resp.Project.Metadata.Keywords,
			Language: resp.Project.Metadata.Language,
			Isbn:     resp.Project.Metadata.Isbn,
		},
	}

	for _, genre := range resp.Project.Genres {
		responseItem.Genres = append(responseItem.Genres, projectGenreItemResponse{
			Id:   genre.Id,
			Name: genre.Name,
		})
	}

	if len(resp.Project.Chapters) == 0 {
		responseItem.Chapters = []chapterResponse{}
	}
	for _, chapter := range resp.Project.Chapters {
		responseItem.Chapters = append(responseItem.Chapters, chapterResponse{
			Id:      chapter.Id,
			Name:    chapter.Name,
			Content: chapter.Content,
			Order:   chapter.Order,
			Version: chapter.Version,
		})
	}

	if len(resp.Project.Characters) == 0 {
		responseItem.Characters = []characterResponse{}
	}
	for _, character := range resp.Project.Characters {
		responseItem.Characters = append(responseItem.Characters, characterResponse{
			Id:          character.Id,
			Name:        character.Name,
			Description: character.Description,
			Attributes:  character.Attributes,
		})
	}

	if len(resp.Project.Locations) == 0 {
		responseItem.Locations = []locationResponse{}
	}
	for _, location := range resp.Project.Locations {
		responseItem.Locations = append(responseItem.Locations, locationResponse{
			Id:          location.Id,
			Name:        location.Name,
			Description: location.Description,
		})
	}

	ctx.JSON(http.StatusOK, responseItem)
}
