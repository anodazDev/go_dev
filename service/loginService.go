package service

type LoginService interface {
	LoginUser(email string, password string) bool
}

type ProfileService interface {
	ProfileUser() (string, string, string)
}

type loginInformation struct {
	email    string
	password string
}

type userInformation struct {
	email    string
	password string
	username string
	tel      string
}

type userProfileInformation struct {
	email    string
	username string
	tel      string
}

func StaticLoginService() LoginService {
	userprofile := UserProfile()
	return &loginInformation{
		email:    userprofile.email,
		password: userprofile.password,
	}
}

func UserProfile() *userInformation {
	return &userInformation{
		email:    "miraistorm@gmail.com",
		password: "password",
		username: "Nick",
		tel:      "0800000000",
	}

}

func (info *loginInformation) LoginUser(email string, password string) bool {
	return info.email == email && info.password == password
}

func ProfileUser() (string, string, string) {
	userprofile := UserProfile()
	return userprofile.email, userprofile.username, userprofile.tel
}
