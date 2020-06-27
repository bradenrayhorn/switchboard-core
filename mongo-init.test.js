db.createUser(
    {
        user: "switchboard",
        pwd: "password",
        roles: [
            {
                role: "readWrite",
                db: "switchboard"
            }
        ]
    }
);
