package main

const HELP_MESSAGE = "Чтобы получить дневник, необходимо внести свои данные с помощью команды /createNewUser. После чего вы сможете получить свой дневник командой /getDiary."
const CREATE_USER_MESSAGE = `Введите, пожалуйста, ваши данные в формате
								Имя
								Группа
								Пульс в спокойном состоянии
								Пульс после нагрузок
								День недели, в который присылать вам напоминание о том, что необходимо распечатать дневник
								Час, в который это делать
							Например, ваша запись может выглядеть следующим образом:
								Иван Иванов
								ИУ7-41
								55
								155
								Понедельник
								14
								`
const REMIND_MESSAGE = "Самое время распечатать дневник."
const USER_NOT_REGISTERED_MESSAGE = "Вы еще не внесли свои данные, так что создать дневник невозможно."
const ERROR_MESSAGE = "Данные некорректны, попробуйте еще раз"
const SUCCESS_USER_CREATION = "Пользователь успешно создан"
const NUMBER_OF_FIELDS_IN_CREATE_USER_MESSAGE = 6