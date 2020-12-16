export default {
    base: "https://api.xutiancheng.me",
    testbase: "https://localhost:8080",
    handlers: {
        users: "/v1/users",
        myuser: "/v1/users/me",
        myuserAvatar: "/v1/users/me/avatar",
        sessions: "/v1/sessions",
        sessionsMine: "/v1/sessions/mine",
        resetPasscode: "/v1/resetcodes",
        passwords: "/v1/passwords/",
        upload: "/v1/upload/"
    }
}