import random

num_groups = 1000
num_people = 100 * 1000
entry_per_person = 10


def do():
    # (group, person, created_at)
    rows = []
    for i in range(num_people):
        for j in range(entry_per_person):
            rows.append((i % num_groups, i, random.randint(0, 10000000)))

    return rows


def insert_statement(rows, table):
    # print(rows)
    vals = ",".join([f"('{row[0]}', '{row[1]}', '{row[2]}')" for row in rows])
    return f"INSERT INTO {table} (`group`, `person`, `created_at`) VALUES {vals};"


rows = do()


batch_size = 1000

with open("insert.sql", "w") as f:
    f.truncate()
    f.write('truncate table test.test1;')

    for i in range(0, len(rows), batch_size):
        rs = rows[i : i + batch_size]
        query = insert_statement(rows, "test.test1")
        f.write(query)

    f.close()
