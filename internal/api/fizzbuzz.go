package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/cruffinoni/fizzbuzz/internal/database"
	"github.com/cruffinoni/fizzbuzz/internal/utils"
	"github.com/gin-gonic/gin"
)

type PlayFizzBuzzBody struct {
	Number1  int    `json:"number1"`
	Number2  int    `json:"number2"`
	Replace1 string `json:"replace1"`
	Replace2 string `json:"replace2"`
	Limit    int    `json:"limit"`
}

type PlayFizzBuzzResponse struct {
	Result string `json:"result"`
}

func formatFizzBuzzFromBody(req *PlayFizzBuzzBody) *PlayFizzBuzzResponse {
	var result string
	for i := 1; i <= req.Limit; i++ {
		if i%req.Number1 == 0 && i%req.Number2 == 0 {
			result += req.Replace1 + req.Replace2
		} else if i%req.Number1 == 0 {
			result += req.Replace1
		} else if i%req.Number2 == 0 {
			result += req.Replace2
		} else {
			result += strconv.Itoa(i)
		}
		result += ","
	}
	return &PlayFizzBuzzResponse{Result: result[:len(result)-1]}
}

func (r *Routes) PlayFizzBuzz(ctx *gin.Context) {
	var fizzBuzzBody PlayFizzBuzzBody
	if err := ctx.ShouldBindJSON(&fizzBuzzBody); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewBadRequestBuilder(err.Error()))
		return
	}
	if fizzBuzzBody.Limit < 1 {
		ctx.JSON(http.StatusBadRequest, utils.NewBadRequestBuilder("limit must be greater than 0"))
		return
	}
	if fizzBuzzBody.Number1 < 1 || fizzBuzzBody.Number2 < 1 {
		ctx.JSON(http.StatusBadRequest, utils.NewBadRequestBuilder("numbers must be greater than 0"))
		return
	}
	if fizzBuzzBody.Replace1 == "" || fizzBuzzBody.Replace2 == "" {
		ctx.JSON(http.StatusBadRequest, utils.NewBadRequestBuilder("replacements cannot be empty"))
		return
	}

	if err := r.db.AddRequest(&database.FizzBuzzRequest{
		Int1:  int64(fizzBuzzBody.Number1),
		Int2:  int64(fizzBuzzBody.Number2),
		Limit: int64(fizzBuzzBody.Limit),
		Str1:  fizzBuzzBody.Replace1,
		Str2:  fizzBuzzBody.Replace2,
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewInternalServerErrorBuilder(err))
		return
	}
	ctx.JSON(http.StatusOK, formatFizzBuzzFromBody(&fizzBuzzBody))
}

func (r *Routes) GetMostUsedRequest(ctx *gin.Context) {
	mostUsed, err := r.db.GetMostUsedRequest()
	if err != nil {
		if errors.Is(err, database.ErrNoRequest) {
			ctx.JSON(http.StatusNotFound, utils.NewStatusOKBuilder("not enough data"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.NewInternalServerErrorBuilder(err))
		return
	}
	ctx.JSON(http.StatusOK, mostUsed)
}
