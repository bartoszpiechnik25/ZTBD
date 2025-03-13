FROM python:3.13-slim-bookworm

WORKDIR /app

ARG BEGIN_YEAR
ARG END_YEAR

ENV BEGIN_YEAR=${BEGIN_YEAR}
ENV END_YEAR=${END_YEAR}

COPY requirements.txt ./

RUN pip install -r requirements.txt

COPY parser.py parser.py

CMD [ "python", "parser.py" ]
