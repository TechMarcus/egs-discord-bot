FROM python:3.13-alpine

WORKDIR /app

COPY . /app

RUN apk add --no-cache gcc musl-dev libffi-dev

RUN pip install --upgrade pip && pip install -r requirements.txt
RUN rm -rf /app/games_info.json

EXPOSE 8080

CMD ["python","-u","main.py"]
