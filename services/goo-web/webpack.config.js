var path = require('path');
var webpack = require('webpack');
var ExtractTextPlugin = require('extract-text-webpack-plugin');
var precss       = require('precss');
var autoprefixer = require('autoprefixer');
var OptimizeCssAssetsPlugin = require('optimize-css-assets-webpack-plugin');

var zlib = require('zlib');
var CompressionPlugin = require("compression-webpack-plugin");

function gzipMaxCompression(buffer, done) {
  return zlib.gzip(buffer, { level: 9 }, done)
}

module.exports = {
  entry: [
    './src/index'
  ],
  output: {
    path: path.join(__dirname, 'dist'),
    filename: 'bundle.js',
    publicPath: '/static/'
  },
  plugins: [
    new webpack.DefinePlugin({
      'process.env.NODE_ENV': JSON.stringify(process.env.NODE_ENV || 'production')
    }),
    new webpack.NoErrorsPlugin(),
    new ExtractTextPlugin('app.css'),
    new webpack.DefinePlugin({
      'process.env.SERVER': `"${process.env.SERVER}"`
    }),
    new OptimizeCssAssetsPlugin({
     assetNameRegExp: /\.css$/g,
     cssProcessor: require('cssnano'),
     cssProcessorOptions: { discardComments: {removeAll: true } },
     canPrint: true
    }),
    new webpack.optimize.UglifyJsPlugin(),
    new webpack.optimize.OccurenceOrderPlugin(),
    new webpack.optimize.DedupePlugin(),
		new CompressionPlugin({
    	algorithm: gzipMaxCompression,
      regExp: /\.(js|html)$/,
      minRatio: 0
    })
  ],
  module: {
    loaders: [
      {
        test: /\.jsx?$/,
        loaders: ['babel'],
        include: path.join(__dirname, 'src')
      },
      { test: /\.css$/, loader: ExtractTextPlugin.extract('css-loader!postcss-loader') },
      { test: /\.(sass|scss)$/, loader: ExtractTextPlugin.extract('css-loader!postcss-loader!sass-loader') },
			{
      	test: /\.html$/,
        name: "mandrillTemplates",
        loader: 'raw!html-minify'
      },
    ]
  },
  postcss: function () {
    return [precss, autoprefixer];
  },
  resolve: {
    extensions: ['', '.js', '.json']
  }
};
