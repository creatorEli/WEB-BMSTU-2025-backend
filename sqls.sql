SELECT * from armies    

INSERT INTO "armies" ("army_id", "name_army", "description_army", "status_army", "image_army_url", "class_army", "min_plain_speed", "max_plain_speed", "min_mount_speed", "max_mount_speed", "min_forest_speed", "max_forest_speed", "min_river_speed", "max_river_speed", "min_desert_speed", "max_desert_speed") VALUES
(6,	'Егеря',	NULL,	'действует',	'http://127.0.0.1:9000/armies/Yeger.jpg',	'step',	25,	35,	15,	20,	12,	18,	5,	15,	15,	25);
-- добавление данных в таблицу услуг - сперва удали пятую услугу 

SELECT * from travel_times