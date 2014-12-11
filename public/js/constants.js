var constantModule = angular.module('uchiwa.constants', []);

// Settings
constantModule.constant('settings', {
  date: 'yyyy-MM-dd HH:mm:ss',
  hideSilenced: true,
  hideOccurrences: true,
  theme: 'default'
});

// Version
constantModule.constant('version', {
  uchiwa: '0.4.0'
});
