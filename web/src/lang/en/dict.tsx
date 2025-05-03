export const dict = {
	common: {
		search: {
			posts: 'Введите название или описание',
			people: 'Введите имя, должность или описанию',
		},
		tabs: {
			posts: 'Топ',
			network: 'Люди',
			collaborations: 'Создать',
			profile: 'Профиль',
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
			save: 'Сохранить',
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
			edit: {
				general: {
					title: 'Расскажите о себе',
					description: 'Коротко о себе, чем занимаетесь и чем увлекаетесь.',
					firstName: 'Имя',
					lastName: 'Фамилия',
					jobTitle: 'Должность',
				},
				description: {
					title: 'О себе',
					description: '',
					placeholder: 'Например: 32 года, серийный предприниматель и директор по продукту с опытом в архитектуре, дизайне, маркетинге и разработке технологий.',
				},
				location: {
					description: 'Укажите ваш город или регион, чтобы другим было проще понять, где вы находитесь.',
					title: 'Где вы живёте?',
				},
				interests: {
					description: 'Расскажите, в каких проектах или темах вам было бы интересно участвовать.',
					title: 'Что вам интересно?',
				},
				badges: {
					description: 'Выберите качества или умения, которые лучше всего вас характеризуют.',
					title: 'Как бы вы себя описали?',
				},
				image: {
					description: 'Добавьте фотографию, чтобы сделать профиль более живым и привлекательным.',
					title: 'Загрузите своё фото',
				},
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
		},
		collaborations: {
			edit: {
				general: {
					description: 'Помогите другим лучше понять вашу задачу',
					title: 'Опишите сотрудничество',
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
