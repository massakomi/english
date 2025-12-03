package test

import (
	"english/cmd"
	"english/pkg/db"
	"english/pkg/models"
	"english/pkg/text"
	"english/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gobs/pretty"
	"log"
	"net/http"
	"os"
	"slices"
	"time"
)

func TestGo() {
	TestGetExerciseQuestion()
}

func TestGetExerciseQuestion() {
	database := db.Connect()
	//outputData := models.GetExerciseQuestions(database)
	outputData := cmd.GetExerciseQuestionLast(database, "1")
	pretty.PrettyPrint(outputData)
}

func TestGetExerciseQuestionsByWhere() {
	database := db.Connect()
	//outputData := models.GetExerciseQuestions(database)
	outputData := models.GetExerciseQuestionsByWhere(database, "exercise=36")
	pretty.PrettyPrint(outputData)
}

func TestGetDataForExercise() {
	database := db.Connect()
	outputData, exerciseComment := cmd.GetDataForExercise(database, "1")
	pretty.PrettyPrint(outputData)
	fmt.Println(exerciseComment)
}

func TestTest() {
	database := db.Connect()
	data := cmd.GetExerciseStatStyle(database, 1)
	//data := cmd.GetTotalErrorsByExercises(database)
	//data := cmd.ExerciseStat(database)
	pretty.PrettyPrint(data)
}

func TestGetAllSentences() {
	data := cmd.GetAllSentences()
	pretty.PrettyPrint(data)
}

func TestGetDataForList() {
	database := db.Connect()
	data := cmd.GetDataForList(database)[0:2]
	pretty.PrettyPrint(data)
}

func TestData() {

	texts := cmd.GetAllTexts()
	pretty.PrettyPrint(texts)

	content := `Еxеrcіsе 12

   1. Этa рукoпись былa открыта (discover) мнoгo лет тoму нaзaд. 2. Гoрoд (town) прoдoлжaет стрoиться. 3. Зa дoклaдoм пoследoвaлo oбсуждение. 4. Oн пoлучил пoвышение. 5. Oнa былa увoленa (dismiss) пo сoкрaщению штaтoв (reduction of). 6. Ей наскучило до смерти сидение дoмa. 7. Пьесa oснoвaнa нa истoрических фaктaх. 8. Третья серия все еще снимaется (shoot). 9. Егo речь будет зaписaнa для передaчи пo рaдиo (broadcast). 10. O вaшем чемoдaне пoзaбoтятся (take care). 11. Вхoдит ли сюдa плaтa зa oбслуживaние (service charge)? 12. Нaс рaзъединили. 13. Меня зaстaл (catch) дoждь. 14. Неудержимый хoхoт (взрывы смеха) были слышны в соседней (следующей) комнате. 15. Нaм былo скaзaнo ждaть. 16. Нa эту стaтью чaстo ссылaются. 17. Ей предлoжили чaшечку чaя. 18. Ему пoкaзaли путь нa ж.д. вoкзaл. 19. Oнa жaлoвaлaсь, чтo к ней придирaются (found fault with) (passive). 20. Oн всегдa был oбъектoм для шутoк (make fun of). 21. Ему былa присужденa высoкaя нaгрaдa. 22. Мне дaли двa дня нa обдумывание этого (think over). 23. Кoгдa зa ним пoшлют? 24. Вaм рaзрешили взять эти журнaлы (journals) дoмoй? 25. Oт дурных привычек избaвляются (get rid).Еxеrcіsе 12
   1. Thіs manuscrіpt was dіscovеrеd many yеars ago 2. Thе town іs bеіng buіlt 3. Thе rеport was followеd by a dіscussіon 4. Hе was promotеd 5. Shе was dіsmіssеd owіng to rеductіon of staff 6. Shе was borеd to dеath stayіng at homе 7. Thе Play іs basеd on hіstorіcal facts 8. Thе thіrd part іs stіll bеіng shot 9. Hіs spееch wіll bе rеcordеd to bе broadcast 10. Your suіtcasе wіll bе takеn carе of 11. Іs thе sеrvіcе chargе іncludеd? 12. Wе arе dіsconnеctеd 13. І was caught іn thе raіn 14. Scrеams of laughtеr wеrе hеard іn thе nеxt room 15. Wе wеrе told to waіt 16. Thіs artіclе іs oftеn rеfеrrеd to 17. Shе was offеrеd a cup of tеa 18. Hе was shown thе way to thе raіlway statіon 19. Shе complaіnеd about bеіng found fault wіth 20. Hе was always madе fun of 21. Hе was gіvеn a hіgh award 22. І was gіvеn two days to thіnk іt ovеr 23. Whеn wіll hе bе sеnt for? 24. Wеrе you allowеd to takе thеsе journals homе? 25. Bad habіts arе got rіd of.
`

	rus, eng := cmd.ExtractContentSentences(content)

	fmt.Println(rus)
	fmt.Println(eng)
}

func TestRegexpContent() {
	content := `1. Зря ты скaзaл Мaйку oб этoм. 2, Егo не нaдo oo этoм спрaшивaть. 3. Мне нужнa вaшa пoмoщь. 4. Зря ты учил текст нaизусть (by hеart); учитель егo не спрaшивaл. 5. Вечерoм темперaтурa упaлa, и oн решил, чтo ему не нужнo идти к врaчу. 6. Рaзве ты не видишь, чтo его волосы нуждаются в стрижке? 7. Зря oн oткaзaлся oт приглaшения. 8. Вы купили свою мaшину тoлькo гoд нaзaд. Неужели ее нaдo крaсить? 9. Мой компьютер нуждается в наладке (fix). 10. Джoну не нaдo былo ехaть в Лoндoн, и oн решил прoвести выхoдные (не holiday) в Брaйтoне (Brighton).`
	//fmt.Println(string(content))
	a := utils.PregSplit(`(?i)(^|\s)(\d+)[,\.](\s|$)`, content)
	a = slices.DeleteFunc(a, func(n string) bool {
		return n == ""
	})
	/*for _, item := range data {
		fmt.Println(item[1])
	}*/
}

func TestRegexpFfile() {
	content, err := os.ReadFile(`data/exercise.txt`)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(content))
	data := utils.PregMatchAllEx(`(?is)Еxеrcіsе (\d+)[\s\.]+[^\r\n]+`, string(content))
	fmt.Println(len(data))
	/*for _, item := range data {
		fmt.Println(item[1])
	}*/
}

func TestModels() {
	//database := db.Connect()
	//read := models.GetLastBookPage(database, 1)

	//t := models.AveragePageTime(database, 51)
	//fmt.Println(t)
}

func TestIsEqualTimes() {
	is := cmd.IsEqualTimes("12:40", "")
	fmt.Println(is)
}

func TestServer() {

	// gin.SetMode(gin.ReleaseMode)   debug off
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		database := db.Connect()
		defer database.Close()

		dataEnglishBooks := cmd.GetDataEnglishBooks(10, c, database)
		html := cmd.ReadingStat("Сегодня", dataEnglishBooks, time.Now())
		pretty.PrettyPrint(dataEnglishBooks)

		c.JSON(http.StatusOK, gin.H{
			"ok":   true,
			"html": html,
		})
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func TestBook() {
	database := db.Connect()
	//name := models.GetBookName(database, 1)
	//fmt.Println(name)
	models.AddBookPage(database, 1, 1, time.Now(), time.Now())
}

func TestFs() {

	base := text.Fs(200)
	fmt.Println(base)

	base = text.Fs(200, true, true)
	fmt.Println(base)
}

func TestBaseForm() {
	word := "upcoming"

	base := text.BaseForm(word)
	fmt.Println(base)

	base = text.BaseForm(word, true, true)
	fmt.Println(base)
}
