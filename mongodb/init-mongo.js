db = db.getSiblingDB('ztbd');
db.createUser({
	user: "tester",
	pwd: "passwd",
	roles: [{ role: "readWrite", db: "ztbd" }]
});
db.createCollection('test_collection');
print('Database initialized: ztbd with user.');
