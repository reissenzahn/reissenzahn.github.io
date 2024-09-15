
```py
import pandas as pd
import numpy as np

# Series

# mean
# s.mean()

# max
# s.max()

# read series from single-column csv
# pd.read_csv('data.csv').squeeze()

# replace NaN values
# s.fillna(4)
```

```py
# Series

## constructor

# from iterable
s = Series([1, 2, 3])

## attributes

# index
s.index

# values as a numpy.ndarray
s.values

# data type
s.dtype

# shape tuple (len,)
s.shape

# number of bytes
s.nbytes

# length
s.size

# memory usage
s.memory_usage()
s.memory_usage(deep=True)

# contains NaN?
s.hasnans

# is empty?
s.empty

# series name
s.name


## Conversion

# cast to dtype
s.astype(np.int32)

# convert to list
s.to_list()

# convert to numpy.ndarray
s.to_numpy()

# copy
s.copy()
s.copy(deep=True)

# convert columns to best possible dtypes
s.convert_dtypes()


## indexing

# access element by label
s.at['C']

# access element by integer position
s.iat[3]

# access range by labels
s.loc['C':'E']

# access range by integer position
s.iloc[3:5]


## binary operator functions
s.mul(3)
s.add(s2)
s.round(decimals=2)


## transformation
s.apply(lambda x: x * 2)


## aggregation
```


Function application, GroupBy & window

Series.agg([func, axis])
Aggregate using one or more operations over the specified axis.

Series.transform(func[, axis])
Call func on self producing a Series with the same axis shape as self.

Series.map(arg[, na_action])

Map values of Series according to an input mapping or function.

Series.groupby([by, axis, level, as_index, ...])

Group Series using a mapper or by a Series of columns.

Series.rolling(window[, min_periods, ...])

Provide rolling window calculations.

Series.expanding([min_periods, axis, method])

Provide expanding window calculations.

Series.ewm([com, span, halflife, alpha, ...])

Provide exponentially weighted (EW) calculations.

Series.pipe(func, *args, **kwargs)

Apply chainable functions that expect Series or DataFrames.


```py
## descriptive stats

# count non-NaN values
s.count()

# descriptive statistics
s.describe()

# summary statistics
s.sum()
s.max()
s.mean()
s.std()
```



Series.nlargest([n, keep])
Return the largest n elements.

Series.nsmallest([n, keep])
Return the smallest n elements.

Series.pct_change([periods, fill_method, ...])
Fractional change between the current and a prior element.

Series.quantile([q, interpolation])
Return value at the given quantile.

Series.rank([axis, method, numeric_only, ...])
Compute numerical data ranks (1 through n) along axis.

Series.unique()
Return unique values of Series object.

Series.nunique([dropna])
Return number of unique elements in the object.


Series.value_counts([normalize, sort, ...])
Return a Series containing counts of unique values.




Reindexing / selection / label manipulation

Series.align(other[, join, axis, level, ...])
Align two objects on their axes with the specified join method.

Series.drop([labels, axis, index, columns, ...])
Return Series with specified index labels removed.

Series.droplevel(level[, axis])

Return Series/DataFrame with requested index / column level(s) removed.
Series.drop_duplicates(*[, keep, inplace, ...])

Return Series with duplicate values removed.
Series.duplicated([keep])

Indicate duplicate Series values.

Series.equals(other)

Test whether two objects contain the same elements.

Series.first(offset)

(DEPRECATED) Select initial periods of time series data based on a date offset.

Series.head([n])

Return the first n rows.

Series.idxmax([axis, skipna])

Return the row label of the maximum value.

Series.idxmin([axis, skipna])

Return the row label of the minimum value.

Series.isin(values)

Whether elements in Series are contained in values.

Series.last(offset)

(DEPRECATED) Select final periods of time series data based on a date offset.

Series.reindex([index, axis, method, copy, ...])

Conform Series to new index with optional filling logic.

Series.reindex_like(other[, method, copy, ...])

Return an object with matching indices as other object.

Series.rename([index, axis, copy, inplace, ...])

Alter Series index labels or name.

Series.rename_axis([mapper, index, axis, ...])

Set the name of the axis for the index or columns.

Series.reset_index([level, drop, name, ...])

Generate a new DataFrame or Series with the index reset.

Series.sample([n, frac, replace, weights, ...])

Return a random sample of items from an axis of object.

Series.set_axis(labels, *[, axis, copy])

Assign desired index to given axis.

Series.take(indices[, axis])

Return the elements in the given positional indices along an axis.

Series.tail([n])

Return the last n rows.

Series.truncate([before, after, axis, copy])

Truncate a Series or DataFrame before and after some index value.

Series.where(cond[, other, inplace, axis, level])

Replace values where the condition is False.

Series.mask(cond[, other, inplace, axis, level])

Replace values where the condition is True.

Series.add_prefix(prefix[, axis])

Prefix labels with string prefix.

Series.add_suffix(suffix[, axis])

Suffix labels with string suffix.

Series.filter([items, like, regex, axis])

Subset the dataframe rows or columns according to the specified index labels.

Missing data handling
Series.backfill(*[, axis, inplace, limit, ...])

(DEPRECATED) Fill NA/NaN values by using the next valid observation to fill the gap.

Series.bfill(*[, axis, inplace, limit, ...])

Fill NA/NaN values by using the next valid observation to fill the gap.

Series.dropna(*[, axis, inplace, how, ...])

Return a new Series with missing values removed.

Series.ffill(*[, axis, inplace, limit, ...])

Fill NA/NaN values by propagating the last valid observation to next valid.

Series.fillna([value, method, axis, ...])

Fill NA/NaN values using the specified method.

Series.interpolate([method, axis, limit, ...])

Fill NaN values using an interpolation method.

Series.isna()

Detect missing values.

Series.isnull()

Series.isnull is an alias for Series.isna.

Series.notna()

Detect existing (non-missing) values.

Series.notnull()

Series.notnull is an alias for Series.notna.

Series.pad(*[, axis, inplace, limit, downcast])

(DEPRECATED) Fill NA/NaN values by propagating the last valid observation to next valid.

Series.replace([to_replace, value, inplace, ...])

Replace values given in to_replace with value.

Reshaping, sorting
Series.argsort([axis, kind, order, stable])

Return the integer indices that would sort the Series values.

Series.argmin([axis, skipna])

Return int position of the smallest value in the Series.

Series.argmax([axis, skipna])

Return int position of the largest value in the Series.

Series.reorder_levels(order)

Rearrange index levels using input order.

Series.sort_values(*[, axis, ascending, ...])

Sort by the values.

Series.sort_index(*[, axis, level, ...])

Sort Series by index labels.

Series.swaplevel([i, j, copy])

Swap levels i and j in a MultiIndex.

Series.unstack([level, fill_value, sort])

Unstack, also known as pivot, Series with MultiIndex to produce DataFrame.

Series.explode([ignore_index])

Transform each element of a list-like to a row.

Series.searchsorted(value[, side, sorter])

Find indices where elements should be inserted to maintain order.

Series.ravel([order])

(DEPRECATED) Return the flattened underlying data as an ndarray or ExtensionArray.

Series.repeat(repeats[, axis])

Repeat elements of a Series.

Series.squeeze([axis])

Squeeze 1 dimensional axis objects into scalars.

Series.view([dtype])

(DEPRECATED) Create a new view of the Series.

Combining / comparing / joining / merging
Series.compare(other[, align_axis, ...])

Compare to another Series and show the differences.

Series.update(other)

Modify Series in place using values from passed Series.

Time Series-related
Series.asfreq(freq[, method, how, ...])

Convert time series to specified frequency.

Series.asof(where[, subset])

Return the last row(s) without any NaNs before where.

Series.shift([periods, freq, axis, ...])

Shift index by desired number of periods with an optional time freq.

Series.first_valid_index()

Return index for first non-NA value or None, if no non-NA value is found.

Series.last_valid_index()

Return index for last non-NA value or None, if no non-NA value is found.

Series.resample(rule[, axis, closed, label, ...])

Resample time-series data.

Series.tz_convert(tz[, axis, level, copy])

Convert tz-aware axis to target time zone.

Series.tz_localize(tz[, axis, level, copy, ...])

Localize tz-naive index of a Series or DataFrame to target time zone.

Series.at_time(time[, asof, axis])

Select values at particular time of day (e.g., 9:30AM).

Series.between_time(start_time, end_time[, ...])

Select values between particular times of the day (e.g., 9:00-9:30 AM).

Accessors
pandas provides dtype-specific methods under various accessors. These are separate namespaces within Series that only apply to specific data types.

Series.str

alias of StringMethods

Series.cat

alias of CategoricalAccessor

Series.dt

alias of CombinedDatetimelikeProperties

Series.sparse

alias of SparseAccessor

DataFrame.sparse

alias of SparseFrameAccessor

Index.str

alias of StringMethods

Data Type

Accessor

Datetime, Timedelta, Period

dt

String

str

Categorical

cat

Sparse

sparse

Datetimelike properties
Series.dt can be used to access the values of the series as datetimelike and return several properties. These can be accessed like Series.dt.<property>.

Datetime properties
Series.dt.date

Returns numpy array of python datetime.date objects.

Series.dt.time

Returns numpy array of datetime.time objects.

Series.dt.timetz

Returns numpy array of datetime.time objects with timezones.

Series.dt.year

The year of the datetime.

Series.dt.month

The month as January=1, December=12.

Series.dt.day

The day of the datetime.

Series.dt.hour

The hours of the datetime.

Series.dt.minute

The minutes of the datetime.

Series.dt.second

The seconds of the datetime.

Series.dt.microsecond

The microseconds of the datetime.

Series.dt.nanosecond

The nanoseconds of the datetime.

Series.dt.dayofweek

The day of the week with Monday=0, Sunday=6.

Series.dt.day_of_week

The day of the week with Monday=0, Sunday=6.

Series.dt.weekday

The day of the week with Monday=0, Sunday=6.

Series.dt.dayofyear

The ordinal day of the year.

Series.dt.day_of_year

The ordinal day of the year.

Series.dt.days_in_month

The number of days in the month.

Series.dt.quarter

The quarter of the date.

Series.dt.is_month_start

Indicates whether the date is the first day of the month.

Series.dt.is_month_end

Indicates whether the date is the last day of the month.

Series.dt.is_quarter_start

Indicator for whether the date is the first day of a quarter.

Series.dt.is_quarter_end

Indicator for whether the date is the last day of a quarter.

Series.dt.is_year_start

Indicate whether the date is the first day of a year.

Series.dt.is_year_end

Indicate whether the date is the last day of the year.

Series.dt.is_leap_year

Boolean indicator if the date belongs to a leap year.

Series.dt.daysinmonth

The number of days in the month.

Series.dt.days_in_month

The number of days in the month.

Series.dt.tz

Return the timezone.

Series.dt.freq

Return the frequency object for this PeriodArray.

Series.dt.unit

Datetime methods
Series.dt.isocalendar()

Calculate year, week, and day according to the ISO 8601 standard.

Series.dt.to_period(*args, **kwargs)

Cast to PeriodArray/PeriodIndex at a particular frequency.

Series.dt.to_pydatetime()

(DEPRECATED) Return the data as an array of datetime.datetime objects.

Series.dt.tz_localize(*args, **kwargs)

Localize tz-naive Datetime Array/Index to tz-aware Datetime Array/Index.

Series.dt.tz_convert(*args, **kwargs)

Convert tz-aware Datetime Array/Index from one time zone to another.

Series.dt.normalize(*args, **kwargs)

Convert times to midnight.

Series.dt.strftime(*args, **kwargs)

Convert to Index using specified date_format.

Series.dt.round(*args, **kwargs)

Perform round operation on the data to the specified freq.

Series.dt.floor(*args, **kwargs)

Perform floor operation on the data to the specified freq.

Series.dt.ceil(*args, **kwargs)

Perform ceil operation on the data to the specified freq.

Series.dt.month_name(*args, **kwargs)

Return the month names with specified locale.

Series.dt.day_name(*args, **kwargs)

Return the day names with specified locale.

Series.dt.as_unit(*args, **kwargs)

Period properties
Series.dt.qyear

Series.dt.start_time

Get the Timestamp for the start of the period.

Series.dt.end_time

Get the Timestamp for the end of the period.

Timedelta properties
Series.dt.days

Number of days for each element.

Series.dt.seconds

Number of seconds (>= 0 and less than 1 day) for each element.

Series.dt.microseconds

Number of microseconds (>= 0 and less than 1 second) for each element.

Series.dt.nanoseconds

Number of nanoseconds (>= 0 and less than 1 microsecond) for each element.

Series.dt.components

Return a Dataframe of the components of the Timedeltas.

Series.dt.unit

Timedelta methods
Series.dt.to_pytimedelta()

Return an array of native datetime.timedelta objects.

Series.dt.total_seconds(*args, **kwargs)

Return total duration of each element expressed in seconds.

Series.dt.as_unit(*args, **kwargs)

```py

```

Series.str can be used to access the values of the series as strings and apply several methods to it. These can be accessed like Series.str.<function/property>.

Series.str.capitalize()

Convert strings in the Series/Index to be capitalized.

Series.str.casefold()

Convert strings in the Series/Index to be casefolded.

Series.str.cat([others, sep, na_rep, join])

Concatenate strings in the Series/Index with given separator.

Series.str.center(width[, fillchar])

Pad left and right side of strings in the Series/Index.

Series.str.contains(pat[, case, flags, na, ...])

Test if pattern or regex is contained within a string of a Series or Index.

Series.str.count(pat[, flags])

Count occurrences of pattern in each string of the Series/Index.

Series.str.decode(encoding[, errors])

Decode character string in the Series/Index using indicated encoding.

Series.str.encode(encoding[, errors])

Encode character string in the Series/Index using indicated encoding.

Series.str.endswith(pat[, na])

Test if the end of each string element matches a pattern.

Series.str.extract(pat[, flags, expand])

Extract capture groups in the regex pat as columns in a DataFrame.

Series.str.extractall(pat[, flags])

Extract capture groups in the regex pat as columns in DataFrame.

Series.str.find(sub[, start, end])

Return lowest indexes in each strings in the Series/Index.

Series.str.findall(pat[, flags])

Find all occurrences of pattern or regular expression in the Series/Index.

Series.str.fullmatch(pat[, case, flags, na])

Determine if each string entirely matches a regular expression.

Series.str.get(i)

Extract element from each component at specified position or with specified key.

Series.str.index(sub[, start, end])

Return lowest indexes in each string in Series/Index.

Series.str.join(sep)

Join lists contained as elements in the Series/Index with passed delimiter.

Series.str.len()

Compute the length of each element in the Series/Index.

Series.str.ljust(width[, fillchar])

Pad right side of strings in the Series/Index.

Series.str.lower()

Convert strings in the Series/Index to lowercase.

Series.str.lstrip([to_strip])

Remove leading characters.

Series.str.match(pat[, case, flags, na])

Determine if each string starts with a match of a regular expression.

Series.str.normalize(form)

Return the Unicode normal form for the strings in the Series/Index.

Series.str.pad(width[, side, fillchar])

Pad strings in the Series/Index up to width.

Series.str.partition([sep, expand])

Split the string at the first occurrence of sep.

Series.str.removeprefix(prefix)

Remove a prefix from an object series.

Series.str.removesuffix(suffix)

Remove a suffix from an object series.

Series.str.repeat(repeats)

Duplicate each string in the Series or Index.

Series.str.replace(pat, repl[, n, case, ...])

Replace each occurrence of pattern/regex in the Series/Index.

Series.str.rfind(sub[, start, end])

Return highest indexes in each strings in the Series/Index.

Series.str.rindex(sub[, start, end])

Return highest indexes in each string in Series/Index.

Series.str.rjust(width[, fillchar])

Pad left side of strings in the Series/Index.

Series.str.rpartition([sep, expand])

Split the string at the last occurrence of sep.

Series.str.rstrip([to_strip])

Remove trailing characters.

Series.str.slice([start, stop, step])

Slice substrings from each element in the Series or Index.

Series.str.slice_replace([start, stop, repl])

Replace a positional slice of a string with another value.

Series.str.split([pat, n, expand, regex])

Split strings around given separator/delimiter.

Series.str.rsplit([pat, n, expand])

Split strings around given separator/delimiter.

Series.str.startswith(pat[, na])

Test if the start of each string element matches a pattern.

Series.str.strip([to_strip])

Remove leading and trailing characters.

Series.str.swapcase()

Convert strings in the Series/Index to be swapcased.

Series.str.title()

Convert strings in the Series/Index to titlecase.

Series.str.translate(table)

Map all characters in the string through the given mapping table.

Series.str.upper()

Convert strings in the Series/Index to uppercase.

Series.str.wrap(width, **kwargs)

Wrap strings in Series/Index at specified line width.

Series.str.zfill(width)

Pad strings in the Series/Index by prepending '0' characters.

Series.str.isalnum()

Check whether all characters in each string are alphanumeric.

Series.str.isalpha()

Check whether all characters in each string are alphabetic.

Series.str.isdigit()

Check whether all characters in each string are digits.

Series.str.isspace()

Check whether all characters in each string are whitespace.

Series.str.islower()

Check whether all characters in each string are lowercase.

Series.str.isupper()

Check whether all characters in each string are uppercase.

Series.str.istitle()

Check whether all characters in each string are titlecase.

Series.str.isnumeric()

Check whether all characters in each string are numeric.

Series.str.isdecimal()

Check whether all characters in each string are decimal.

Series.str.get_dummies([sep])

Return DataFrame of dummy/indicator variables for Series.



List accessor
Arrow list-dtype specific methods and attributes are provided under the Series.list accessor.

Series.list.flatten()

Flatten list values.

Series.list.len()

Return the length of each list in the Series.

Series.list.__getitem__(key)

Index or slice lists in the Series.

Struct accessor
Arrow struct-dtype specific methods and attributes are provided under the Series.struct accessor.

Series.struct.dtypes

Return the dtype object of each child field of the struct.

Series.struct.field(name_or_index)

Extract a child field of a struct as a Series.

Series.struct.explode()

Extract all child fields of a struct as a DataFrame.


Plotting
Series.plot is both a callable method and a namespace attribute for specific plotting methods of the form Series.plot.<kind>.

Series.plot([kind, ax, figsize, ....])

Series plotting accessor and method

Series.plot.area([x, y, stacked])

Draw a stacked area plot.

Series.plot.bar([x, y])

Vertical bar plot.

Series.plot.barh([x, y])

Make a horizontal bar plot.

Series.plot.box([by])

Make a box plot of the DataFrame columns.

Series.plot.density([bw_method, ind])

Generate Kernel Density Estimate plot using Gaussian kernels.

Series.plot.hist([by, bins])

Draw one histogram of the DataFrame's columns.

Series.plot.kde([bw_method, ind])

Generate Kernel Density Estimate plot using Gaussian kernels.

Series.plot.line([x, y])

Plot Series or DataFrame as lines.

Series.plot.pie(**kwargs)

Generate a pie plot.

Series.hist([by, ax, grid, xlabelsize, ...])

Draw histogram of the input series using matplotlib.


```py
## Serialization

# csv
s.to_csv('data.csv')

# DataFrame
s.to_frame()

# json
s.to_json('data.json')

# markdown
print(s.to_markdown(index=False))
```

