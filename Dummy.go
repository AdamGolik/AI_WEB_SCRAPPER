package main

func Dummy() {
	var body struct {
		Name            string
		Lastname        string
		Email           string
		Password        string
		PasswordConfirm string
	}
	println(body.Email)
	println(body.Lastname)
	println(body.Name)
	println(body.Password)
	println(body.PasswordConfirm)
	println("Dummy function called")
}

