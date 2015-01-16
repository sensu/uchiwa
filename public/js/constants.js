var constantModule = angular.module('uchiwa.constants', []);

// Settings
constantModule.constant('settings', {
  date: 'yyyy-MM-dd HH:mm:ss',
  hideSilenced: false,
  hideOccurrences: false,
  theme: 'default'
});

// Version
constantModule.constant('version', {
  uchiwa: '0.4.1'
});
