# a little tool to join CSVs

Use it like this: `joincsv labels.csv content1.csv content2.csv`

Labels from the first CSV are applied as headers for the successive CSVs.
The first row of the first CSV should be your desired new headers in order.
The subsequent rows of the first CSV should be those same headers in the 
columns where those values appear in your content CSVs.

If that doesn't make any sense, see the examples folder. 