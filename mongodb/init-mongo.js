db = db.getSiblingDB('ztbd');
db.createUser({
	user: "tester",
	pwd: "passwd",
	roles: [{ role: "readWrite", db: "ztbd" }]
});
db.createCollection("awards")
print('Database initialized: ztbd with user: tester.');
