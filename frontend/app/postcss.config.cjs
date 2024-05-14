// postcss.config.js
const plugins = {
	'postcss-rem': {
		name: 'toRem'
	},
	autoprefixer: {},
	cssnano: { preset: 'default' }
};

module.exports = {
	plugins
};
