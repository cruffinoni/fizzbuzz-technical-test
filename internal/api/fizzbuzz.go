package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/cruffinoni/fizzbuzz/internal/database"
	"github.com/cruffinoni/fizzbuzz/internal/utils"
	"github.com/gin-gonic/gin"
)

type PlayFizzBuzzBody struct {
	Number1            int64  `json:"number1" binding:"required"`
	Number2            int64  `json:"number2" binding:"required"`
	Replace1           string `json:"replace1" binding:"required"`
	Replace2           string `json:"replace2" binding:"required"`
	Limit              int64  `json:"limit" binding:"required"`
	DisablePerformance bool   `json:"disable_performance"`
}

type PlayFizzBuzzResponse struct {
	Result string `json:"result"`
}

const softLimit = 10000

func formatFizzBuzzFromBody(start int64, req *PlayFizzBuzzBody) *PlayFizzBuzzResponse {
	var result string
	for i := start; i <= req.Limit; i++ {
		result += insertFizzBuzzAlg(i, req) + ","
	}
	if result == "" {
		return &PlayFizzBuzzResponse{Result: ""}
	}
	return &PlayFizzBuzzResponse{Result: result[:len(result)-1]}
}

func basicFormatFizzBuzzFromBody(req *PlayFizzBuzzBody) *PlayFizzBuzzResponse {
	return formatFizzBuzzFromBody(1, req)
}

func insertFizzBuzzAlg(i int64, req *PlayFizzBuzzBody) string {
	if i%req.Number1 == 0 && i%req.Number2 == 0 {
		return req.Replace1 + req.Replace2
	} else if i%req.Number1 == 0 {
		return req.Replace1
	} else if i%req.Number2 == 0 {
		return req.Replace2
	} else {
		return strconv.FormatInt(i, 10)
	}
}

func formatFizzBuzzFromBodyWithPerformance(req *PlayFizzBuzzBody) *PlayFizzBuzzResponse {
	var (
		routines    = req.Limit / softLimit
		remainder   = req.Limit % softLimit
		resultOrder = make([]string, routines+1)
		wg          sync.WaitGroup
		mu          sync.Mutex
	)
	wg.Add(int(routines))
	for i := int64(0); i < routines; i++ {
		go func(routineNb int64) {
			defer wg.Done()
			start := routineNb*softLimit + 1
			limit := (routineNb + 1) * softLimit
			localRes := formatFizzBuzzFromBody(start, &PlayFizzBuzzBody{
				Number1:  req.Number1,
				Number2:  req.Number2,
				Replace1: req.Replace1,
				Replace2: req.Replace2,
				Limit:    limit,
			})
			mu.Lock()
			resultOrder[routineNb] = localRes.Result
			mu.Unlock()
		}(i)
	}
	if remainder > 0 {
		start := routines*softLimit + 1
		localRes := formatFizzBuzzFromBody(start, &PlayFizzBuzzBody{
			Number1:  req.Number1,
			Number2:  req.Number2,
			Replace1: req.Replace1,
			Replace2: req.Replace2,
			Limit:    req.Limit,
		})
		resultOrder[routines] = localRes.Result
	}
	wg.Wait()
	res := strings.Join(resultOrder, ",")
	l := len(res)
	if res[l-1] == ',' {
		res = res[:l-1]
	}
	return &PlayFizzBuzzResponse{Result: res}
}

// PlayFizzBuzz godoc
//
//	@Summary		Play custom FizzBuzz
//	@Description	Returns a customized fizz-buzz list
//	@ID				play-fizzbuzz
//	@Accept			json
//	@Produce		json
//	@Param			b	body		PlayFizzBuzzBody	true	"Required parameters to play fizz-buzz"
//	@Success		200	{object}	PlayFizzBuzzResponse
//	@Router			/play [post]
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
		Int1:  fizzBuzzBody.Number1,
		Int2:  fizzBuzzBody.Number2,
		Limit: fizzBuzzBody.Limit,
		Str1:  fizzBuzzBody.Replace1,
		Str2:  fizzBuzzBody.Replace2,
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewInternalServerErrorBuilder(err))
		return
	}
	if fizzBuzzBody.Limit > softLimit && !fizzBuzzBody.DisablePerformance {
		ctx.JSON(http.StatusOK, formatFizzBuzzFromBodyWithPerformance(&fizzBuzzBody))
	} else {
		ctx.JSON(http.StatusOK, basicFormatFizzBuzzFromBody(&fizzBuzzBody))
	}
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
