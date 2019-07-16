		permissions = append(permissions, bomodel.Permission{
			Method:      "POST",
			Endpoint:    "/account",
			Function:    c.Create,
			Description: "Create New Account",
			Category:    "ACCOUNT",
		})

		permissions = append(permissions, bomodel.Permission{
			Method:      "GET",
			Endpoint:    "/account",
			Function:    c.GetAll,
			Description: "Get All Account",
			Category:    "ACCOUNT",
		})

		permissions = append(permissions, bomodel.Permission{
			Method:      "GET",
			Endpoint:    "/account/:accountID",
			Function:    c.GetOne,
			Description: "Get One Account",
			Category:    "ACCOUNT",
		})

		permissions = append(permissions, bomodel.Permission{
			Method:      "DELETE",
			Endpoint:    "/account/:accountID",
			Function:    c.Delete,
			Description: "Delete a Account",
			Category:    "ACCOUNT",
		})

		permissions = append(permissions, bomodel.Permission{
			Method:      "PUT",
			Endpoint:    "/account/:accountID",
			Function:    c.Update,
			Description: "Update a Account",
			Category:    "ACCOUNT",
		})


	uuidv4Regex := regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}")

	regexColonID := regexp.MustCompile(":[a-zA-Z0-9]+")

	permissionMaps := map[[2]string]string{}

	apiController := router.Group("/api")

	apiController.Use(func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		s := strings.Fields(token)
		if len(s) != 2 {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		if s[0] != "Bearer" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		jwt := tk.ValidateToken("LOGIN", s[1])
		if jwt == nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Set("jwt", jwt)

		data, ok := jwt.ExtendData.(map[string]interface{})
		if !ok {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		userID := data["userId"].(string)
		c.Set("user.id", userID)

		urlx := uuidv4Regex.ReplaceAllString(c.Request.URL.Path, "#")

		permissionCode, ok := permissionMaps[[2]string{c.Request.Method, urlx}]
		if !ok {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		spaceID := c.Param("spaceId")

		userType := data["type"].(string)

		if userType == "ADMIN" {
			return
		}

		sc := bomodel.ServiceContext{"user.id": userID}

		isAdmin := adminService.IsAdmin(sc, bomodel.IsAdminRequest{SpaceID: spaceID})

		if isAdmin {
			return
		}

		isAccessable := userService.IsAccessable(sc, bomodel.IsAccessableRequest{
			MethodEndpoint: permissionCode,
			SpaceID:        spaceID,
		})

		if !isAccessable {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

	})

	for _, permission := range permissions {
		apiController.Handle(permission.Method, permission.Endpoint, permission.Function)
		hashedIDURL := regexColonID.ReplaceAllString(apiController.BasePath()+permission.Endpoint, "#")
		permissionMaps[[2]string{permission.Method, hashedIDURL}] = fmt.Sprintf("%s_%s", permission.Method, hashedIDURL)
	}

	// everyone free to access it
	router.POST("/register", guestController.Register)
	router.POST("/login", guestController.Login)
	router.GET("/activate", guestController.Activate)
	router.POST("/password/forgot/init", guestController.ForgotPasswordInit)
	router.POST("/password/forgot/reset", guestController.ForgotPasswordReset)

	userService.CreateAdminUserIfNotExist(nil)