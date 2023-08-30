package src

import "context"

const (
	APILink1 = "https://run.mocky.io/v3/03d2a7bd-f12f-4275-9e9a-84e41f9c2aae"
	APILink2 = "https://run.mocky.io/v3/aab281fe-3dbb-4d91-a863-a96e6bf083d7"
)

func Build() {
	newUsecase()
}

func GetFromSource() (data []UserData, err error) {
	return uc.GetSampleAPIResourceRedirect(context.Background(), []string{
		APILink1,
		APILink2,
	})
}

func SetAndReplaceToCSV(data []UserData, path string) error {
	return uc.StoreAndReplaceUserDataToCSV(context.Background(), data, path)
}

func SearchFromCSV(tags []string, path string) (data []UserData, err error) {
	return uc.SearchUserWithTags(context.Background(), tags, path)
}
