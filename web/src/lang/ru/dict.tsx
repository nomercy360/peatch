export const dict = {
	common: {
		tabs: {
			network: 'Люди',
			collaborations: 'Новый проект',
			posts: 'Лента',
		},
		search: {
			posts: 'Искать посты',
			people: 'Искать людей',
			noMoreResults: 'Больше нет результатов',
			noResults: 'Ничего не найдено',
		},
		buttons: {
			generateRandomAvatar: 'Создать случайный аватар',
			expressInterest: 'Проявить интерес',
			startCollaboration: 'Начать сотрудничество',
			next: 'Далее',
			chooseAndSave: 'Выбрать и сохранить',
			publish: 'Опубликовать',
			done: 'Готово',
			edit: 'Редактировать',
			backToProfile: 'Назад к профилю',
			clear: 'Очистить',
			skip: 'Пропустить',
		},
		textarea: {
			maxLength: 'максимум {{ maxLength }} символов',
		},
	},
	pages: {
		notFound: {
			title: '404: Страница не найдена',
		},
		users: {
			verificationStatusDenied: 'Мы скрыли ваш профиль. Попробуйте сделать его более личным и настоящим.',
			shareURLText: 'Посмотри профиль {{name}} на Peatch! 🌟',
			edit: {
				general: {
					title: 'Расскажите о себе',
					fullName: 'Ваше полное имя',
					jobTitle: 'Тайтл работы',
				},
				description: {
					title: 'О себе',
					placeholder: 'Например: 32 года, серийный предприниматель и директор по продукту с опытом в архитектуре, дизайне, маркетинге и разработке технологий.',
				},
				location: {
					title: 'Где вы живёте?',
				},
				interests: {
					title: 'Что вам интересно?',
				},
				badges: {
					title: 'Как бы вы себя описали?',
				},
				image: {
					title: 'Загрузите своё фото',
				},
			},
			fillProfilePopup: {
				title: 'Заполните профиль',
				description: 'Заполните профиль всего за 5 минут, чтобы расширить возможности для общения и сотрудничества.',
				action: 'Заполнить профиль',
			},
			collaborate: {
				title: 'Проявить интерес',
				description: 'Напишите сообщение, чтобы начать сотрудничество',
			},
			activity: {
				title: 'Уведомления',
			},
			availableFor: 'Готов к',
			sayHi: 'Привет',
			saidHi: 'Запрос отправлен',
			shareProfile: 'Поделиться',
			followSuccess: 'Мы отправили уведомление о том, что вам понравился их профиль',
			followError: 'Не удалось отправить уведомление',
			botBlocked: 'Пользователю нельзя отправить уведомление через бота',
			messageUser: 'Написать в Telegram',
		},
		collaborations: {
			edit: {
				general: {
					description: 'Помогите другим лучше понять вашу задачу',
					title: 'Опишите проект',
					titlePlaceholder: 'Ищу дизайнера продукта',
					descriptionPlaceholder: 'Например: Ищу дизайнера для участия в некоммерционном хакатоне',
					checkboxPlaceholder: 'Оплачивается ли эта возможность?',
				},
				location: {
					title: 'Есть ли предпочтения по локации?',
					description: 'Укажите место, которое лучше всего подходит для сотрудничества',
				},
				interests: {
					title: 'Выберите тему',
					description: 'Это поможет нам рекомендовать вашу инициативу другим',
					chooseOne: 'выберите один',
					selectedCount: '{{count}} из 10',
					searchPlaceholder: 'Поиск возможностей для сотрудничества',
				},
				badges: {
					title: 'Кого вы ищете?',
					description: 'Выберите теги, которые лучше всего описывают вашу задачу',
					searchPlaceholder: 'Поиск по тегам',
				},
				createBadge: {
					title: 'Создание {{ name }}',
					description: 'Это поможет нам рекомендовать вас другим пользователям',
				},
			},
		},
	},
	components: {
		actionDonePopup: {
			success: 'Успешно',
			callToAction: 'Продолжить',
		},
	},
} as const
