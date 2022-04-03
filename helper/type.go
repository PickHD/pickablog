package helper

import "github.com/PickHD/pickablog/model"

// BuildMetaData will building a meta data responses based on param
func BuildMetaData(page int, size int, order string, totalData int, totalPage int) *model.Metadata {
	return &model.Metadata{
		Page: page,
		Size: size,
		Order: order,
		TotalData: totalData,
		TotalPage: totalPage,
	}
}