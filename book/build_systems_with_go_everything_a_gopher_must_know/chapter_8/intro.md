# CHAPTER 8: ENCODINGS

We can find popular formats such as CSV, JSON, XML or YAML that are used to represent data from byte level to a human-readable formart.
Other formats such as **base64** or **PEM** data serizlization are oriented to facilitate machine to machine interaction.

In **Go**, the package `encoding` offers a set of subpackages that facilitate the conversion from Go types to these  formats and the other way around.

## 8.1 CSV

**Comma Separated Values (CSV) is a widely used format to represent tabular data. Go provides CSV read/write operators in the package `encoding/csv`. Every line from a CSV is defined as a CSV record and every record contains the items separated by commas.

In order to read CSV data we use a `Reader` type that converts CSV lines into CSV records. j
