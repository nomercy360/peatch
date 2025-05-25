export const dict = {
  common: {
    tabs: {
      network: 'People',
      collaborations: 'New project',
      posts: 'Feed',
    },
    search: {
      posts: 'Enter title or description',
      people: 'Enter name, job title or description',
      noMoreResults: 'No more results',
      noResults: 'Nothing found',
    },
    buttons: {
      generateRandomAvatar: 'Generate random avatar',
      expressInterest: 'Show interest',
      startCollaboration: 'Start collaboration',
      next: 'Next',
      chooseAndSave: 'Select and save',
      publish: 'Publish',
      done: 'Done',
      edit: 'Edit',
      backToProfile: 'Back to profile',
      clear: 'Clear',
      save: 'Save',
      skip: 'Skip',
      cancel: 'Cancel',
      update: 'Update',
      add: 'Add',
    },
    textarea: {
      maxLength: 'maximum {{ maxLength }} characters',
    },
  },
  pages: {
    notFound: {
      title: '404: Page not found',
    },
    users: {
      verificationStatusDenied:
        'We have hidden your profile. Try to make it more personal and authentic.',
      shareURLText: "Check out {{name}}'s profile on Peatch! 🌟",
      edit: {
        general: {
          fullName: 'Full name',
          title: 'Tell about yourself',
          description:
            'Briefly about yourself, what you do and what interests you.',
          firstName: 'Full name',
          jobTitle: 'Job title',
        },
        description: {
          title: 'About me',
          description:
            'Tell about yourself, what you do and what interests you.',
          placeholder:
            'For example: 32 years old, serial entrepreneur and product director with experience in architecture, design, marketing and technology development.',
        },
        location: {
          description:
            'Specify your city or region to help others understand where you are located.',
          title: 'Where do you live?',
        },
        interests: {
          description:
            'Tell us what projects or topics you would be interested in participating in.',
          title: 'What interests you?',
        },
        badges: {
          description: 'Choose qualities or skills that best describe you.',
          title: 'How would you describe yourself?',
        },
        image: {
          description:
            'Add a photo to make your profile more lively and attractive.',
          title: 'Upload your photo',
        },
      },
      fillProfilePopup: {
        title: 'Set up your profile',
        description:
          'Complete your profile in just 5 minutes to enhance your networking and be able to collaborate with others.',
        action: 'Set up profile',
      },
      collaborate: {
        title: 'Show interest',
        description: 'Write a message to start collaboration',
      },
      activity: {
        title: 'Notifications',
      },
      availableFor: 'Available for',
      sayHi: 'Say Hi',
      saidHi: 'Request sent',
      shareProfile: 'Share',
      followSuccess:
        'We have sent a notification that you appreciate their profile',
      followError: 'Failed to send notification',
      botBlocked: 'User cannot be notified via bot',
      messageUser: 'Message in Telegram',
      publishSuccess: 'Your profile is now visible to others',
      publishError: 'Failed to publish profile',
      profileIncomplete: 'Please complete your profile before publishing',
      profileBlocked: 'Blocked profiles cannot be published',
      links: 'Links',
      linksUpdated: 'Links updated successfully',
      linksUpdateError: 'Failed to update links',
      invalidUrl: 'Please enter a valid URL',
      addLink: 'Add link',
      linkUrlPlaceholder: 'Enter URL',
      linkLabelPlaceholder: 'Enter link title',
      fetchingMetadata: 'Fetching link info...',
      showMore: 'Show more',
      showLess: 'Show less',
      editLink: 'Edit Link',
      linkUrl: 'URL',
      linkLabel: 'Title',
    },
    collaborations: {
      edit: {
        general: {
          description: 'Help others better understand your task',
          title: 'Describe the project',
          titlePlaceholder: 'Looking for a product designer',
          descriptionPlaceholder:
            'For example: Looking for a designer to participate in a non-profit hackathon',
          checkboxPlaceholder: 'Is this opportunity paid?',
        },
        location: {
          title: 'Do you have location preferences?',
          description: 'Specify the place that best suits the collaboration',
        },
        interests: {
          title: 'Choose a topic',
          description: 'This will help us recommend your initiative to others',
          chooseOne: 'choose one',
          selectedCount: '{{count}} of 10',
          searchPlaceholder: 'Search for collaboration opportunities',
        },
        badges: {
          title: 'Who are you looking for?',
          description: 'Choose tags that best describe your task',
          searchPlaceholder: 'Search by tags',
        },
        createBadge: {
          title: 'Creating {{ name }}',
          description: 'This will help us recommend you to other users',
        },
      },
    },
  },
  components: {
    actionDonePopup: {
      success: 'Success',
      callToAction: 'Continue',
    },
  },
} as const
