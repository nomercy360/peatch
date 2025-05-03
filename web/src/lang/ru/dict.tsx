export const dict = {
	common: {
		search: {
			posts: 'Искать посты',
			people: 'Искать людей',
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
					firstName: 'Имя',
					lastName: 'Фамилия',
					jobTitle: 'Должность',
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
