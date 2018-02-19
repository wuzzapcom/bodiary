package main

const helpMessage = "Чтобы получить дневник, необходимо внести свои данные с помощью команды /createNewUser. После чего вы сможете получить свой дневник командой /getDiary."
const createUserMessage = `Введите, пожалуйста, ваши данные в формате
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
							Если вы не хотите получать уведомления, то не указывайте последние два пункта.
								`
const remindMessage = "Самое время распечатать дневник."
const userNotRegisteredMessage = "Вы еще не внесли свои данные, так что создать дневник невозможно."
const errorMessage = "Данные некорректны, попробуйте еще раз"
const successUserCreation = "Пользователь успешно создан"
const minimalNumberOfFieldsInCreateUserMessage = 4
const pathToUsersFoldersFlag = "pathToUsersFolders"
const pathToUsersFoldersHelp = "Папка, в которой будут храниться метаданные пользователей."
const pathToAuthFileFlag = "pathToAuthFile"
const pathToAuthFileHelp = "Путь к файлу, в котором находится ключ телеграм-бота."
const pathToHTMLTemplateFlag = "pathToHtmlTemplate"
const pathToHTMLTempleteHelp = "Путь к файлу, который является go-совместимым HTML-темплейтом."
const pathToGeneratedFilesFlag = "pathToGeneratedFiles"
const pathToGeneratedFilesHelp = "Путь к папке, где будут храниться сгенерированные HTML-файлы дневников."
const helpCommand = "help"
const createNewUserCommand = "create_new_user"
const getDiaryCommand = "get_diary"
