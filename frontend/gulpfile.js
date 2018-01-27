const gulp = require('gulp');
const less = require('gulp-less');
const babel = require('gulp-babel');
const rename = require('gulp-rename');
const uglify = require('gulp-uglify');
const webpack = require('gulp-webpack');
const cleanCss = require('gulp-clean-css');

gulp.task('frontend-styles', () => {
  return gulp.src('styles/main.less')
    .pipe(less({paths: ['node_modules/bootstrap/less']}))
    .pipe(cleanCss())
    .pipe(rename('bundled.css'))
    .pipe(gulp.dest('public'));
});

gulp.task('frontend-libs', () => {
  return gulp.src('lib/**/*.js')
    .pipe(babel({
      presets: ['es2015'],
    }))
    .pipe(gulp.dest('public/scripts/lib'));
});

gulp.task('frontend-components', () => {
  return gulp.src('components/**/*.jsx')
    .pipe(babel({
      presets: ['es2015', 'stage-0'],
      plugins: ['transform-react-jsx']
    }))
    .pipe(gulp.dest('public/scripts/components'));
});

gulp.task('frontend-scripts', ['frontend-libs', 'frontend-components'], () => {
  return gulp.src('public/scripts/components/index.js')
    .pipe(webpack({
      output: {
        filename: 'bundled.js'
      }
    }))
    .pipe(gulp.dest('public'));
});

gulp.task('watch', ['default'], () => {
  gulp.watch(['components/**/*.jsx', 'lib/**/*.js'], ['frontend-scripts']);
  gulp.watch(['styles/*.less'], ['frontend-styles']);
});

gulp.task('default', ['frontend-styles', 'frontend-scripts']);
